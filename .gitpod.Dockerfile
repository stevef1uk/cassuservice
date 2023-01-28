FROM gitpod/workspace-full

RUN sudo apt-get update \
    && brew tap go-swagger/go-swagger \
    && brew install yq \
    && brew install go-swagger && sudo ln -s /home/linuxbrew/.linuxbrew/bin/swagger /usr/local/bin/swagger

