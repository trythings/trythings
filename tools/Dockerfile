FROM golang:1.6.4
ENV CODE /go/src/github.com/trythings/trythings

# Install less.
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends less \
	&& rm -rf /var/lib/apt/lists/*

# Install docker.
RUN apt-get update \
	&& apt-get install --yes --no-install-recommends \
		apt-transport-https \
		ca-certificates \
		gnupg2 \
	&& rm -rf /var/lib/apt/lists/* \
	&& apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D \
	&& echo "deb https://apt.dockerproject.org/repo debian-jessie main" > /etc/apt/sources.list.d/docker.list \
	&& apt-get update \
	&& apt-get install --yes --no-install-recommends docker-engine \
	&& rm -rf /var/lib/apt/lists/* \
	&& service docker start

# Install git-lfs.
RUN curl --fail --silent --show-error --location https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash \
	&& apt-get install --yes --no-install-recommends git-lfs \
	&& rm -rf /var/lib/apt/lists/*

# Install fish.
RUN echo "deb http://download.opensuse.org/repositories/shells:/fish:/release:/2/Debian_8.0/ /" > /etc/apt/sources.list.d/fish.list \
	&& wget --quiet http://download.opensuse.org/repositories/shells:fish:release:2/Debian_8.0/Release.key \
	&& apt-key add - < Release.key \
	&& rm Release.key \
	&& apt-get update \
	&& apt-get install --yes --no-install-recommends fish \
	&& rm -rf /var/lib/apt/lists/*

# Copy the code.
COPY . "$CODE"

# Install fisherman.
RUN mkdir -p "$HOME"/.config/fish/functions \
	&& ln -s "$CODE"/tools/home/.config/fish/functions/fisher.fish "$HOME"/.config/fish/functions/fisher.fish

# Install bass.
RUN fish --command "fisher edc/bass"

# Install fish.config.
RUN mkdir -p "$HOME"/.config/fish \
	&& ln -s "$CODE"/tools/home/.config/fish/config.fish "$HOME"/.config/fish/config.fish

CMD ["fish"]
