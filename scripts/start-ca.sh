#!/bin/bash

fabric-ca-server start --cfg.identities.allowremove --cfg.affiliations.allowremove -b admin:adminpw -d
