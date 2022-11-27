// Copyright (C) 2022 Tiago Ferreira - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"
	"fmt"

	"github.com/reviewpad/cookbook/codehost"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type Size struct {
	targetEntity handler.TargetEntity
	codehost     codehost.Codehost
	collector    collector.Collector
	log          *logrus.Entry
}

func NewSizeRecipe(targetEntity handler.TargetEntity, codehost codehost.Codehost, collector collector.Collector) (*Size, error) {
	return &Size{
		targetEntity,
		codehost,
		collector,
		logrus.WithField("user", fmt.Sprintf("%s/%s", targetEntity.Owner, targetEntity.Repo)),
	}, nil
}

func (s *Size) Run(ctx context.Context) error {
	owner := s.targetEntity.Owner
	repo := s.targetEntity.Repo
	number := s.targetEntity.Number

	s.log.Info("running size recipe")

	err := s.collector.Collect("run recipe", s.collectionData())
	if err != nil {
		s.log.WithError(err).Error("error collecting data")
	}

	if err := createSizeLabels(ctx, owner, repo, s.codehost); err != nil {
		s.log.WithError(err).Error("error creating size labels")
		return err
	}

	prSizeData, err := s.codehost.GetPRSizeData(ctx, owner, repo, number)
	if err != nil {
		s.log.WithError(err).Error("error getting pr size data")
		return err
	}

	labelsToAdd := make([]string, 0)
	labelsToRemove := make([]string, 0)

	if prSizeData.Changes <= 100 {
		labelsToAdd = append(labelsToAdd, "small")
		labelsToRemove = append(labelsToRemove, "medium", "large")
	} else if prSizeData.Changes >= 100 && prSizeData.Changes <= 500 {
		labelsToAdd = append(labelsToAdd, "medium")
		labelsToRemove = append(labelsToRemove, "small", "large")
	} else {
		labelsToAdd = append(labelsToAdd, "large")
		labelsToRemove = append(labelsToRemove, "small", "medium")
	}

	s.log.WithField("labels", labelsToAdd).Info("adding labels")

	s.log.WithField("labels", labelsToRemove).Info("removing labels")

	labels := append(prSizeData.Labels, labelsToAdd...)

	for _, labelToRemove := range labelsToRemove {
		index := slices.Index(labels, labelToRemove)

		if index != -1 {
			labels = slices.Delete(labels, index, index+1)
		}
	}

	s.log.WithField("labels", labels).Info("final labels")

	err = s.codehost.SetLabels(ctx, owner, repo, number, labels)
	if err != nil {
		s.log.WithError(err).Error("error setting pr labels")
		return err
	}

	return nil
}

func (s *Size) collectionData() map[string]interface{} {
	return map[string]interface{}{
		"recipe_name": "size",
		"owner":       s.targetEntity.Owner,
		"repo":        s.targetEntity.Repo,
		"kind":        s.targetEntity.Kind,
		"number":      s.targetEntity.Number,
		"distinct_id": s.targetEntity.Owner,
	}
}

func createSizeLabels(ctx context.Context, owner, repo string, ch codehost.Codehost) error {
	labels := []codehost.Label{
		{
			Name:        "small",
			Color:       "219ebc",
			Description: "Pull request change is between 0 - 100 changes",
		},
		{
			Name:        "medium",
			Color:       "faedcd",
			Description: "Pull request change is between 101 - 500 changes",
		},
		{
			Name:        "large",
			Color:       "e76f51",
			Description: "Pull request change is more than 501+ changes",
		},
	}

	return ch.CreateLabels(ctx, owner, repo, labels)
}
