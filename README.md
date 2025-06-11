# GitOps-Driven Chaos Engineering Framework

**Orchestrate, schedule, and report** Kubernetes chaos experiments via CRDs—fully managed through GitOps.

###  Core Capabilities
- **Declarative Experiments**: Network latency, pod-kill, disk‐fill defined as LitmusChaos CRDs.  
- **GitOps Delivery**: ArgoCD auto-syncs any new experiment in `chaos-playbook/`.  
- **Quality Gates**: GitHub Actions runs `kubeval` & `chaos-lint` on PRs.  
- **Controlled Blast Radius**: Terraform spins up an Azure DevTest Lab for safe testing.  
- **Observability**: Prometheus scrapes chaos metrics; Grafana dashboards visualize SLO impacts.

###  Tech Stack & Flow
1. **Terraform** ➔ Azure DevTest Lab  
2. **Chaos CRDs** ➔ Commit new YAML to `chaos-playbook/`  
3. **GitHub Actions** ➔ Lint & validate on PR  
4. **ArgoCD** ➔ Sync CRDs into target cluster  
5. **Prometheus & Grafana** ➔ Scrape metrics & render “Chaos Metrics” dashboard  

###  Prerequisites
- Kubernetes cluster with LitmusChaos installed  
- ArgoCD deployed to cluster  
- Azure subscription + Contributor role  
- GitHub repo with secrets for Terraform

###  Quickstart
```bash
git clone https://github.com/your-org/chaos-engineering-framework.git
cd terraform && terraform init && terraform apply
cd ../azure-devtest-lab && terraform init && terraform apply
# Open a PR adding your experiment under chaos-playbook/
# GitHub Actions will lint → ArgoCD will deploy → observe in Grafana

