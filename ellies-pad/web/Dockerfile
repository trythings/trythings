FROM node:7.2.1
ENV CODE /go/src/github.com/trythings/trythings

# Copy the code.
COPY . "$CODE"
WORKDIR "$CODE"

# Install dependencies.
RUN npm install

# Run the webapp.
WORKDIR "$CODE"/ellies-pad/web
EXPOSE 3000
ENTRYPOINT ["npm", "start"]
