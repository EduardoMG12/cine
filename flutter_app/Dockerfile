FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    openjdk-8-jdk \
    wget \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME="/usr/lib/jvm/java-8-openjdk-amd64"
ENV FLUTTER_HOME="/opt/flutter"
ENV PATH="$FLUTTER_HOME/bin:$PATH"

RUN git clone https://github.com/flutter/flutter.git -b stable --depth 1 $FLUTTER_HOME
RUN flutter doctor --android-licenses || true
RUN flutter config --enable-web
RUN flutter doctor

WORKDIR /app

COPY pubspec.yaml ./
RUN flutter pub get

COPY . .

EXPOSE 3000

CMD ["flutter", "run", "--hot", "--web-hostname", "0.0.0.0", "--web-port", "3000"]
