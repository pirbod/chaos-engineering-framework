# Ansible Operations

These playbooks are intentionally lightweight. They support bootstrap checks and repeatable operational maintenance without requiring access to a live cluster for the local demo.

## Bootstrap

```bash
ansible-playbook -i ansible/inventory.example.ini ansible/bootstrap-platform.yml
```

The bootstrap playbook checks for common platform CLIs and prints the next Helm/Grafana steps.

## Secret Rotation

```bash
ansible-playbook -i ansible/inventory.example.ini ansible/rotate-secrets.yml
```

The rotation playbook is a safe checklist scaffold for Vault, GitLab and Datadog credentials. Real rotations should be backed by Vault policies, protected CI variables and an audited GitOps change.
