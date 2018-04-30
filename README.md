Helmet
======

Helmet (*Dutch: Helm*) is a project that tries to bring the awesome templating capabilities from [Helm](https://github.com/kubernetes/helm) to Kubernetes without the need of a Tiller. 

Currently this comes at a few downsides like the lack of state which inhibits any rollbacks. It also fully relies on Kubernetes to apply the needed changes and will probably not notice any failed deoployments. 

*Helmet is currently a 1 person side project, so don't expect this to be production ready soon*

## Why not Helm?
When implementing Helm in production I noticed some issues with it while it has awesome capabilities. For my own use I basicly only needed the templating, value merging and a few other small things hich should have been easy to replicate.  
(I also like a coding challenge to understand systems better)

## Thank you
This code is heavily inspired by Helm. Thank you to everybody who contributed to this project!