package controllers

import "github.com/Nedopro2022/wao-estimator/pkg/estimator"

func (r *EstimatorReconciler) GetEstimators() *estimator.Estimators {
	return r.estimators
}
