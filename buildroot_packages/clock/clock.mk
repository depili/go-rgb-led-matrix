#############################################################
#
# Led matrix clock
#
#############################################################
CLOCK_VERSION = 1.0
CLOCK_SITE_METHOD = local
CLOCK_SITE = $(TOPDIR)/package/clock

define CLOCK_BUILD_CMDS
     rm -fr $(@D)/src/github.com/depili/go-rgb-led-matrix
     mkdir -p $(@D)/src/github.com/depili/go-rgb-led-matrix
     git clone https://github.com/depili/go-rgb-led-matrix.git $(@D)/src/github.com/depili/go-rgb-led-matrix
     cd $(@D)/src/github.com/depili/go-rgb-led-matrix && git checkout -b udp origin/udp
     GOPATH=$(@D) GOARCH=arm GOOS=linux GOARM=7 go get github.com/depili/go-rgb-led-matrix/clock
endef

define CLOCK_INSTALL_TARGET_CMDS
     mkdir -p $(TARGET_DIR)/root/rpi-matrix
     $(INSTALL) -D -m 0755 $(@D)/bin/linux_arm/clock $(TARGET_DIR)/root/rpi-matrix/clock

     # init scripts
     $(INSTALL) -D -m 0755 $(@D)/S03copy_clock_files $(TARGET_DIR)/etc/init.d/S03copy_clock_files
     $(INSTALL) -D -m 0755 $(@D)/S03usb_serial $(TARGET_DIR)/etc/init.d/S03usb_serial
     $(INSTALL) -D -m 0755 $(@D)/S99clock $(TARGET_DIR)/etc/init.d/S99clock
     
     $(INSTALL) -D -m 0755 $(@D)/clock_pokemon.sh $(TARGET_DIR)/root/clock_pokemon.sh
     echo "/root/rpi-matrix/clock $(BR2_PACKAGE_CLOCK_CONFIG)" > $(@D)/clock_cmd.sh
     $(INSTALL) -D -m 0755 $(@D)/clock_cmd.sh $(TARGET_DIR)/root/rpi-matrix/clock_cmd.sh
endef

$(eval $(generic-package))
