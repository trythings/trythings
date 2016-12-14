FROM golang:1.6.4
ENV CODE /go/src/github.com/trythings/trythings

# Install gcloud.
RUN apt-get update \
    && apt-get install --yes --no-install-recommends \
        apt-transport-https \
    && rm -rf /var/lib/apt/lists/* \
    && echo "deb https://packages.cloud.google.com/apt cloud-sdk-jessie main" >> /etc/apt/sources.list.d/google-cloud-sdk.list \
    && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
    && apt-get update \
    && apt-get install --yes --no-install-recommends google-cloud-sdk \
    && rm -rf /var/lib/apt/lists/*

# Install goapp.
RUN apt-get update \
    && apt-get install --yes --no-install-recommends unzip \
    && rm -rf /var/lib/apt/lists/*

ENV GO_APPENGINE_VERSION 1.9.48
ENV GO_APPENGINE_DOWNLOAD_URL https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-"$GO_APPENGINE_VERSION".zip
ENV GO_APPENGINE_DOWNLOAD_SHA1 d8f94d4f8c03dec058b3235b3f0248b90475cdf9

RUN curl --fail --silent --show-error --location "$GO_APPENGINE_DOWNLOAD_URL" --output go_appengine.zip \
	&& echo "$GO_APPENGINE_DOWNLOAD_SHA1" go_appengine.zip | sha1sum --check - \
    && unzip -q go_appengine.zip -d /usr/local \
    && rm go_appengine.zip

ENV PATH /usr/local/go_appengine:"$PATH"

# Copy the code.
COPY . "$CODE"

# Run the api.
EXPOSE 8000 8080
WORKDIR "$CODE"/ellies-pad/api/main
RUN echo "opt_in: false" >> "$HOME"/.appcfg_nag \
    && echo "timestamp: 0.0" >> "$HOME"/.appcfg_nag

ENTRYPOINT [ \
    "dev_appserver.py", \
    "--admin_host", "0.0.0.0", \
    "--env_var", "GOOGLE_APPLICATION_CREDENTIALS=../../tools/cmd/deploy/key.json", \
    "--host", "0.0.0.0", \
    "." \
]
