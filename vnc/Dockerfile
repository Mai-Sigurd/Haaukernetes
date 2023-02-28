FROM kalilinux/kali-rolling:latest

# Install dependencies
RUN apt update && apt install kali-linux-core tightvncserver dbus-x11 -y
RUN apt install kali-desktop-xfce -y
RUN apt install -y kali-grant-root && dpkg-reconfigure kali-grant-root

RUN groupadd -r kali && useradd -r -s /bin/bash -g kali kali
RUN echo "kali:kali" | chpasswd
RUN usermod -aG sudo kali

#these lines were switched
USER kali
WORKDIR /home/kali

# Set up VNC
RUN mkdir ~/.vnc
RUN echo "kali" | vncpasswd -f > ~/.vnc/passwd
RUN chmod 600 ~/.vnc/passwd

ENV DISPLAY :1
ENV USER kali
ENV HOME /home/kali

EXPOSE 5901/tcp
#CMD bash -c "rm -f /tmp/.X1-lock /tmp/.X11-unix/X1 && vncserver :1 -geometry 1280x800 -depth 24 && tail -f /dev/null"
CMD bash -c "rm -f /tmp/.X1-lock /tmp/.X11-unix/X1 && vncserver :1 -geometry 1920x1080 -depth 24 && tail -f /dev/null"
