FROM gitpod/workspace-full

RUN sudo apt-get update \
    && brew tap go-swagger/go-swagger \
    && brew install go-swagger 

