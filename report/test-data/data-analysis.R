install.packages("ggplot2")
library(ggplot2)

home_path <- "/Users/theakjeldsmark/Outside OneDrive/Haaukernetes/report/test-data/data/"

# Test 6 challenges and Wireguard one user

# CPU utilization dataset
data <- read.table(paste(home_path, "test6challengesandwireguardoneuser/mem.csv", sep=""), header = TRUE, sep = ",", colClasses = c("numeric", "numeric"), col.names = c("time", "memory"))

# Convert the time column to POSIXct format
data$time <- as.POSIXct(data$time / 1000, origin = "1970-01-01")

# Scatter plot
plot(data$memory ~ data$time, ylab="CPUs", data=data,
     xlab="Time",
     main="Mem Usage Over Time")

vert_line <- as.POSIXct(1683224240000 / 1000, origin = "1970-01-01")
abline(v = vert_line, col = "red", lwd=2)



vert_line <- as.POSIXct(1683224240000/1000, origin = "1970-01-01")

# create ggplot with date and memory data
ggplot(data, aes(x = time, y = memory)) +
  geom_point() +
  geom_vline(xintercept = as.numeric(vert_line), color = "red", linetype = "dashed") +
  labs(x = "Date", y = "Memory") +
  ggtitle("Memory Over Time")