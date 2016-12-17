FROM golang:1.6.4
ENV CODE /go/src/github.com/trythings/trythings

# Start installing from apt-get.
RUN apt-get update

# Install sudo.
# RUN apt-get install --yes --no-install-recommends sudo

# Install less.
RUN apt-get install --yes --no-install-recommends less

# Install docker.
RUN apt-get install --yes --no-install-recommends \
		apt-transport-https \
		ca-certificates \
		gnupg2 \
	&& apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D \
	&& echo "deb https://apt.dockerproject.org/repo debian-jessie main" > /etc/apt/sources.list.d/docker.list \
	&& apt-get update \
	&& apt-get install --yes --no-install-recommends docker-engine \
	&& service docker start

# Install git-lfs.
RUN curl --fail --silent --show-error --location https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash \
	&& apt-get install --yes --no-install-recommends git-lfs

# Install fish.
RUN echo "deb http://download.opensuse.org/repositories/shells:/fish:/release:/2/Debian_8.0/ /" > /etc/apt/sources.list.d/fish.list \
	&& wget --quiet http://download.opensuse.org/repositories/shells:fish:release:2/Debian_8.0/Release.key \
	&& apt-key add - < Release.key \
	&& rm Release.key \
	&& apt-get update \
	&& apt-get install --yes --no-install-recommends fish

# Done installing from apt-get.
RUN rm -rf /var/lib/apt/lists/*

# Install unison.
# ENV UNISON_VERSION 2.48.4
# ENV UNISON_DOWNLOAD_URL https://github.com/bcpierce00/unison/releases/download/"$UNISON_VERSION"/unison-"$UNISON_VERSION"-linux-x86_64-text-static.tar.gz
# ENV UNISON_MD5 7bc592b2acea3070301b580f71e27f9c
# RUN curl --fail --silent --show-error --location "$UNISON_DOWNLOAD_URL" --output unison.tar.gz \
# 	&& echo "$UNISON_MD5" unison.tar.gz | md5sum --check - \
#     && tar -C /usr/local/bin -xzf unison.tar.gz \
# 	&& mv /usr/local/bin/unison-"$UNISON_VERSION"-linux-x86_64-text-static /usr/local/bin/unison \
#     && rm unison.tar.gz

# Copy the code.
COPY . "$CODE"

# Create and run as the user.
ARG USER
RUN useradd --create-home --shell /usr/bin/fish "$USER"
USER "$USER"

# Set git email.
RUN git config --global user.email "$USER"@trythings.com

# Install fisherman.
RUN mkdir -p "$HOME"/.config/fish/functions \
	&& ln -s "$CODE"/tools/home/.config/fish/functions/fisher.fish "$HOME"/.config/fish/functions/fisher.fish

# Install bass.
RUN fish --command "fisher edc/bass"

# Install fish.config.
RUN mkdir -p "$HOME"/.config/fish \
	&& ln -s "$CODE"/tools/home/.config/fish/config.fish "$HOME"/.config/fish/config.fish

WORKDIR "$CODE"
# EXPOSE 5000
# CMD ["unison", "-socket", "5000"]
