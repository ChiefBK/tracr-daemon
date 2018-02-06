# tracr-daemon
The daemon module of Tracr collects and analyzes information for use by tracr-bots.  It is the beating heart of the Tracr  trading platform


## Usage
Tracr-daemon is first inititialized using the Linux systemd service. Before tracr-daemon can be used its service needs to be started manually.

    sudo systemctl start tracr_daemond
    
To stop

    sudo systemctl stop tracr_daemond
    
To restart

    sudo systemctl restart tracr_daemond
    
Once the service is successfully started it stays in a "waiting" state. 

---

The service will remain in this waiting state until the user would like to start collection and analysis. 

To start collection/analysis use command

    tracrd start
    
To stop

    tracrd stop
    
