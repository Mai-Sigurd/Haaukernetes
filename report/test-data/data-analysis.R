# Test 6 challenges and Wireguard one user

# CPU utilization dataset
data <- read.table("data/test6challengesandwireguardoneuser/cpu.csv", header = TRUE, sep = ",", colClasses = c("numeric", "numeric"), col.names = c("time", "memory"))

# convert the time column to POSIXct format
data$time <- as.POSIXct(data$time / 1000, origin = "1970-01-01")

# Scatter plot
plot(data$memory ~ data$time, ylab="CPUs", data=data,
     xlab="Time",
     main="CPU Usage Over Time")

# Memory dataset