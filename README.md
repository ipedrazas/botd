# Botd

[![Build Status](http://drone.ipedrazas.k8s.co.uk/api/badges/ipedrazas/botd/status.svg)](http://drone.ipedrazas.k8s.co.uk/ipedrazas/botd)

Experimenting with convention over configuration in the CI/CD space


## Deploy using Kb8or
         ${SUDO} docker run -i --rm=true \
            -v $(pwd)/k8s:/var/lib/deploy \
            -e KUBECONFIG=/var/lib/deploy/tmp/kb8or/.kube/config \
            quay.io/ukhomeofficedigital/kb8or:v0.4.2 deploy
