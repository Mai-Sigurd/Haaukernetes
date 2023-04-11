FROM kalilinux/kali-rolling:latest

ENV DEBIAN_FRONTEND noninteractive

ENV DESKTOP_ENVIRONMENT=xfce
ENV DESKTOP_PKG=kali-desktop-${DESKTOP_ENVIRONMENT}

ENV REMOTE_ACCESS=rdp

ENV KALI_PACKAGE=core
ENV KALI_PKG=kali-linux-${KALI_PACKAGE}

ENV DISPLAY :1
ENV USER kali
ENV HOME /home/kali

# Install common packages
RUN apt update -q --fix-missing
RUN apt upgrade -y
RUN apt -y install --no-install-recommends sudo wget curl dbus-x11 xinit openssh-server ${DESKTOP_PKG}

# Create a start script
RUN echo "#!/bin/bash" > /startkali.sh
RUN echo "/etc/init.d/ssh start" >> /startkali.sh
RUN chmod 755 /startkali.sh

# Install VNC as a way to keep the container running
RUN apt install tightvncserver dbus-x11 -y
RUN echo "rm -f /tmp/.X1-lock /tmp/.X11-unix/X1 && vncserver :1 -geometry 1920x1080 -depth 24 && tail -f /dev/null" >> /startkali.sh

# Install Kali packages
RUN apt -y install --no-install-recommends ${KALI_PKG}

# Install VNC as a way to keep the container running
RUN apt install tightvncserver dbus-x11 -y
RUN echo "rm -f /tmp/.X1-lock /tmp/.X11-unix/X1 && vncserver :1 -geometry 1920x1080 -depth 24 && tail -f /dev/null" >> /startkali.sh

# Create a non-root Kali user
RUN useradd -m -s /bin/bash -G sudo kali
RUN echo "kali:kali" | chpasswd

# Install and configure xrdp for rdp access
RUN apt -y install --no-install-recommends xorg xorgxrdp xrdp
RUN echo "/etc/init.d/xrdp start" >> /startkali.sh
RUN sed -i s/^port=3389/port=13389/ /etc/xrdp/xrdp.ini
RUN adduser xrdp ssl-cert
RUN echo xfce4-session > /home/kali/.xsession
RUN chmod +x /home/kali/.xsession

# Set up VNC
RUN mkdir ~/.vnc
RUN echo "kali" | vncpasswd -f > ~/.vnc/passwd
RUN chmod 600 ~/.vnc/passwd

# Run start script
RUN echo "/bin/bash" >> /startkali.sh

# Expose ports
EXPOSE 13389
WORKDIR "/root"
ENTRYPOINT ["/bin/bash"]
CMD ["/startkali.sh"]