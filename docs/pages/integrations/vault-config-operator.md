---
layout: enhanced
title: Vault Config Operator
subtitle: Collaborative and visual infrastructure as code for Vault Config Operator
image: /assets/img/integrations/vault-config-operator.svg
permalink: extensibility/integrations/vault-config-operator
docURL: https://docs.meshery.io/extensibility/integrations/vault-config-operator
description: 
category: Provisioning
subcategory: Key Management
registrant: artifacthub
components: 
	-	name: AuthEngineMount
		colorIcon: assets/img/integrations/components/AuthEngineMount-color.svg
		whiteIcon: assets/img/integrations/components/AuthEngineMount-white.svg
		description: 
	-	name: DatabaseSecretEngineConfig
		colorIcon: assets/img/integrations/components/DatabaseSecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/DatabaseSecretEngineConfig-white.svg
		description: 
	-	name: DatabaseSecretEngineRole
		colorIcon: assets/img/integrations/components/DatabaseSecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/DatabaseSecretEngineRole-white.svg
		description: 
	-	name: DatabaseSecretEngineStaticRole
		colorIcon: assets/img/integrations/components/DatabaseSecretEngineStaticRole-color.svg
		whiteIcon: assets/img/integrations/components/DatabaseSecretEngineStaticRole-white.svg
		description: 
	-	name: GitHubSecretEngineConfig
		colorIcon: assets/img/integrations/components/GitHubSecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/GitHubSecretEngineConfig-white.svg
		description: 
	-	name: GitHubSecretEngineRole
		colorIcon: assets/img/integrations/components/GitHubSecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/GitHubSecretEngineRole-white.svg
		description: 
	-	name: JWTOIDCAuthEngineConfig
		colorIcon: assets/img/integrations/components/JWTOIDCAuthEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/JWTOIDCAuthEngineConfig-white.svg
		description: 
	-	name: JWTOIDCAuthEngineRole
		colorIcon: assets/img/integrations/components/JWTOIDCAuthEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/JWTOIDCAuthEngineRole-white.svg
		description: 
	-	name: KubernetesAuthEngineConfig
		colorIcon: assets/img/integrations/components/KubernetesAuthEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/KubernetesAuthEngineConfig-white.svg
		description: 
	-	name: KubernetesAuthEngineRole
		colorIcon: assets/img/integrations/components/KubernetesAuthEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/KubernetesAuthEngineRole-white.svg
		description: 
	-	name: KubernetesSecretEngineConfig
		colorIcon: assets/img/integrations/components/KubernetesSecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/KubernetesSecretEngineConfig-white.svg
		description: 
	-	name: KubernetesSecretEngineRole
		colorIcon: assets/img/integrations/components/KubernetesSecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/KubernetesSecretEngineRole-white.svg
		description: 
	-	name: LDAPAuthEngineConfig
		colorIcon: assets/img/integrations/components/LDAPAuthEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/LDAPAuthEngineConfig-white.svg
		description: 
	-	name: LDAPAuthEngineGroup
		colorIcon: assets/img/integrations/components/LDAPAuthEngineGroup-color.svg
		whiteIcon: assets/img/integrations/components/LDAPAuthEngineGroup-white.svg
		description: 
	-	name: PasswordPolicy
		colorIcon: assets/img/integrations/components/PasswordPolicy-color.svg
		whiteIcon: assets/img/integrations/components/PasswordPolicy-white.svg
		description: 
	-	name: PKISecretEngineConfig
		colorIcon: assets/img/integrations/components/PKISecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/PKISecretEngineConfig-white.svg
		description: 
	-	name: PKISecretEngineRole
		colorIcon: assets/img/integrations/components/PKISecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/PKISecretEngineRole-white.svg
		description: 
	-	name: QuaySecretEngineConfig
		colorIcon: assets/img/integrations/components/QuaySecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/QuaySecretEngineConfig-white.svg
		description: 
	-	name: QuaySecretEngineRole
		colorIcon: assets/img/integrations/components/QuaySecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/QuaySecretEngineRole-white.svg
		description: 
	-	name: QuaySecretEngineStaticRole
		colorIcon: assets/img/integrations/components/QuaySecretEngineStaticRole-color.svg
		whiteIcon: assets/img/integrations/components/QuaySecretEngineStaticRole-white.svg
		description: 
	-	name: RabbitMQSecretEngineConfig
		colorIcon: assets/img/integrations/components/RabbitMQSecretEngineConfig-color.svg
		whiteIcon: assets/img/integrations/components/RabbitMQSecretEngineConfig-white.svg
		description: 
	-	name: RabbitMQSecretEngineRole
		colorIcon: assets/img/integrations/components/RabbitMQSecretEngineRole-color.svg
		whiteIcon: assets/img/integrations/components/RabbitMQSecretEngineRole-white.svg
		description: 
	-	name: RandomSecret
		colorIcon: assets/img/integrations/components/RandomSecret-color.svg
		whiteIcon: assets/img/integrations/components/RandomSecret-white.svg
		description: 
	-	name: SecretEngineMount
		colorIcon: assets/img/integrations/components/SecretEngineMount-color.svg
		whiteIcon: assets/img/integrations/components/SecretEngineMount-white.svg
		description: 
	-	name: VaultSecret
		colorIcon: assets/img/integrations/components/VaultSecret-color.svg
		whiteIcon: assets/img/integrations/components/VaultSecret-white.svg
		description: 
	-	name: Policy
		colorIcon: assets/img/integrations/components/Policy-color.svg
		whiteIcon: assets/img/integrations/components/Policy-white.svg
		description: 
	-	name: GroupAlias
		colorIcon: assets/img/integrations/components/GroupAlias-color.svg
		whiteIcon: assets/img/integrations/components/GroupAlias-white.svg
		description: 
	-	name: Group
		colorIcon: assets/img/integrations/components/Group-color.svg
		whiteIcon: assets/img/integrations/components/Group-white.svg
		description: 
featureList: [
  "Drag-n-drop cloud native infrastructure designer to configure, model, and deploy your workloads.",
  "Invite anyone to review and make changes to your private designs.",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs.
language: en
list: include
---
<p>

</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
