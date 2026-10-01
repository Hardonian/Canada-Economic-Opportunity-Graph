# Release deployment contract

The release workflow publishes and Cosign-signs both API and web images. Its
`image-digests` artifact is the only permitted image input for a production
deployment. Replace the version tags in `k8s/deployment.yaml` (or provide the
equivalent Helm values) with the artifact's full `@sha256:` references before
applying manifests.

The checked-in Kubernetes API deployment uses the WAL-backed persistent mode.
It must remain at one replica because its PVC is `ReadWriteOnce`. For
horizontally scalable read traffic, deploy `STORAGE_MODE=snapshot` and immutable
data releases to each replica; do not share the WAL across API pods.
