FROM node:8.11.3

# Upgrade npm
RUN npm install -g npm

# Install less
RUN npm install -g less

# Install npm-cli-login
RUN npm install -g npm-cli-login

RUN mkdir -p $GOPATH/src/github.com/khulnasoft-lab/

RUN ln -s /khulnasoft-lab $GOPATH/src/github.com/khulnasoft-lab/fastnode

CMD ["bash"]
