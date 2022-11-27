// Copyright (C) 2022 Tiago Ferreira - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"
	"fmt"
	"log"

	"github.com/reviewpad/cookbook/codehost"
	"github.com/reviewpad/reviewpad/v3/collector"
	"github.com/reviewpad/reviewpad/v3/handler"
	"golang.org/x/exp/slices"
)

type Size struct {
	targetEntity handler.TargetEntity
	codehost     codehost.Codehost
	collector    collector.Collector
}

func NewSizeRecipe(targetEntity handler.TargetEntity, codehost codehost.Codehost, collector collector.Collector) (*Size, error) {
	return &Size{
		targetEntity,
		codehost,
		collector,
	}, nil
}

func (s *Size) Run(ctx context.Context) error {
	owner := s.targetEntity.Owner
	repo := s.targetEntity.Repo
	number := s.targetEntity.Number

	err := s.collector.Collect("run recipe", s.collectionData())
	if err != nil {
		log.Println(err)
	}

	if err := createSizeLabels(ctx, owner, repo, s.codehost); err != nil {
		return fmt.Errorf("error creating size labels: %w", err)
	}

	prSizeData, err := s.codehost.GetPRSizeData(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("error getting pr size data: %w", err)
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

	log.Printf("adding labels: %v\n", labelsToAdd)

	log.Printf("removing labels: %v\n", labelsToRemove)

	labels := append(prSizeData.Labels, labelsToAdd...)

	for _, labelToRemove := range labelsToRemove {
		index := slices.Index(labels, labelToRemove)

		if index != -1 {
			labels = slices.Delete(labels, index, index+1)
		}
	}

	log.Printf("final labels to set: %v", labels)

	err = s.codehost.SetLabels(ctx, owner, repo, number, labels)
	if err != nil {
		return fmt.Errorf("error setting pr labels: %w", err)
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
