FROM kalilinux/kali-rolling:latest

ENV DEBIAN_FRONTEND noninteractive

ENV DESKTOP_ENVIRONMENT=xfce
ENV DESKTOP_PKG=kali-desktop-${DESKTOP_ENVIRONMENT}

ENV REMOTE_ACCESS=rdp

ENV KALI_PACKAGE=core
ENV KALI_PKG=kali-linux-${KALI_PACKAGE}

# Install common packages
RUN apt update -q --fix-missing
RUN apt upgrade -y
RUN apt -y install --no-install-recommends sudo wget curl dbus-x11 xinit openssh-server ${DESKTOP_PKG}

# Create a start script
RUN echo "#!/bin/bash" > /startkali.sh
RUN chmod 755 /startkali.sh

# Install Kali packages
RUN apt -y install --no-install-recommends ${KALI_PKG}

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

# Keep the container running
RUN echo "tail -f /dev/null" >> /startkali.sh

# Expose ports
EXPOSE 13389
WORKDIR "/root"
ENTRYPOINT ["/bin/bash"]
CMD ["/startkali.sh"]