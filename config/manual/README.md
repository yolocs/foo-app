```
gcloud iam service-accounts add-iam-policy-binding \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:cshou-playground.svc.id.goog[foo/foo-runner]" \
  foo-runner@cshou-playground.iam.gserviceaccount.com
```