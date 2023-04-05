FROM kalilinux/kali-rolling:latest

ENV DEBIAN_FRONTEND noninteractive
ENV DESKTOP_PKG=kali-desktop-xfce
ENV REMOTE_ACCESS=rdp

ARG UNAME=kali
ARG UPASS=kali
ARG RDP_PORT=3389

# Kali packages to install
ENV KALI_PKG=kali-linux-default

# Install common packages
RUN apt update -q --fix-missing
RUN apt upgrade -y
RUN apt -y install --no-install-recommends sudo wget curl dbus-x11 xinit openssh-server ${DESKTOP_PKG}

# Create a start script
RUN echo "#!/bin/bash" > /startkali.sh
RUN echo "/etc/init.d/ssh start" >> /startkali.sh
RUN chmod 755 /startkali.sh

# Install Kali packages
RUN apt -y install --no-install-recommends ${KALI_PKG}

# Create a non-root Kali user
RUN useradd -m -s /bin/bash -G sudo ${UNAME}
RUN echo "${UNAME}:${UPASS}" | chpasswd

# Install and configure xrdp for rdp access
RUN apt -y install --no-install-recommends xorg xorgxrdp xrdp
RUN echo "/etc/init.d/xrdp start" >> /startkali.sh
RUN sed -i s/^port=3389/port=${RDP_PORT}/ /etc/xrdp/xrdp.ini
RUN adduser xrdp ssl-cert
RUN echo xfce4-session > /home/${UNAME}/.xsession
RUN chmod +x /home/${UNAME}/.xsession

# Run and keep the start script running
RUN echo "/bin/bash" >> /startkali.sh

# Expose ports
EXPOSE ${RDP_PORT}
WORKDIR "/home/kali"
ENTRYPOINT ["/bin/bash"]
CMD ["/startkali.sh"]