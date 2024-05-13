export account_name=tgs-create-site
export namespace=staging

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: $account_name
  namespace: $namespace
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: $account_name
  namespace: $namespace
subjects:
- kind: ServiceAccount
  name: $account_name
  namespace: $namespace
roleRef:
  kind: Role
  name: $account_name
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: $namespace
  name: $account_name
rules:
- apiGroups: [""]
  resources:
  - pods
  - services
  - deployments
  - ingresses
  verbs:
  - get
  - list
  - watch
  - create
  - update
EOF
