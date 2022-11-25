// Copyright (C) 2022 Tiago Ferreira - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package recipes

import (
	"context"

	"github.com/dukex/mixpanel"
	"github.com/reviewpad/cookbook/codehost"
	"github.com/reviewpad/reviewpad/v3/handler"
	"golang.org/x/exp/slices"
)

type Size struct {
	targetEntity handler.TargetEntity
	codehost     codehost.Codehost
	collector    mixpanel.Mixpanel
}

func NewSizeRecipe(targetEntity handler.TargetEntity, codehost codehost.Codehost, collector mixpanel.Mixpanel) (*Size, error) {
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

	if err := createSizeLabels(ctx, owner, repo, s.codehost); err != nil {
		return err
	}

	prSizeData, err := s.codehost.GetPRSizeData(ctx, owner, repo, number)
	if err != nil {
		return err
	}

	labelsToAdd := make([]string, 0)
	labelsToRemove := make([]string, 0)

	if prSizeData.Changes < 100 {
		labelsToAdd = append(labelsToAdd, "small")
		labelsToRemove = append(labelsToRemove, "medium", "large")
	} else if prSizeData.Changes > 100 && prSizeData.Changes < 500 {
		labelsToAdd = append(labelsToAdd, "medium")
		labelsToRemove = append(labelsToRemove, "small", "large")
	} else {
		labelsToAdd = append(labelsToAdd, "large")
		labelsToRemove = append(labelsToRemove, "small", "medium")
	}

	labels := append(prSizeData.Labels, labelsToAdd...)

	for _, labelToRemove := range labelsToRemove {
		index := slices.Index(labels, labelToRemove)

		if index != -1 {
			labels = slices.Delete(labels, index, index+1)
		}
	}

	err = s.codehost.SetLabels(ctx, owner, repo, number, labels)
	if err != nil {
		return err
	}

	return nil
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
