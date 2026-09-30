#!/usr/bin/env bash

set -euo pipefail

chart_path="${CHART_PATH:-deploy/charts/auto-workflow}"
release_name="${RELEASE_NAME:-auto-workflow}"
image_repository="${IMAGE_REPOSITORY:-auto-workflow}"
image_tag="${IMAGE_TAG:-latest}"
ingress_enabled="${INGRESS_ENABLED:-false}"
auto_workflow_env="${AUTO_WORKFLOW_ENV:-devin}"

case "${auto_workflow_env}" in
devin|uat|prod)
	;;
*)
	echo "AUTO_WORKFLOW_ENV must be devin, uat, or prod" >&2
	exit 1
	;;
esac

kube_namespace="${KUBE_NAMESPACE:-auto-workflow-${auto_workflow_env}}"

if [[ "${ingress_enabled}" == "true" && -z "${INGRESS_HOST:-}" ]]; then
	echo "INGRESS_HOST is required when INGRESS_ENABLED=true" >&2
	exit 1
fi

helm lint "${chart_path}"

helm_args=(
	--set-string "image.repository=${image_repository}"
	--set-string "image.tag=${image_tag}"
	--set "ingress.enabled=${ingress_enabled}"
	--set-string "appConfig.environment=${auto_workflow_env}"
)

if [[ "${ingress_enabled}" == "true" ]]; then
	helm_args+=(--set-string "ingress.host=${INGRESS_HOST}")
fi

helm upgrade --install "${release_name}" "${chart_path}" \
	--namespace "${kube_namespace}" \
	--create-namespace \
	--wait \
	--atomic \
	--timeout "${HELM_TIMEOUT:-5m}" \
	"${helm_args[@]}"
