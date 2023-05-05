# install.packages("ggplot2")
library(ggplot2)

home_path <- "/Users/theakjeldsmark/Outside OneDrive/"
test <- "test6challengesandwireguardoneuser"
path <- paste0(home_path, "Haaukernetes/report/test-data/data/one-user/", test, "/")
s <- 180

test_name <- "one Wireguard user with 6 challenges"


# CPU ----------------------------------------------------------------------

# Create an empty data frame with seconds 1 to 180
df <- data.frame(seconds = seq(1, s))

# Create an empty data frame to store max CPU values
max_cpu_df <- data.frame(seconds = numeric(), max_cpu = numeric()) 

# Loop through the files and add the CPU column to the data frame
for (i in 1:5) {
  file_name <- paste0(path, "cpu-", i, ".csv")
  cpu_data <- read.csv(file_name)
  cpu_col <- cpu_data[, 2][1:s] # Extract second column and first 180 rows
  df <- cbind(df, cpu_col)
  
  max_cpu <- max(cpu_col) # Get the max CPU value
  max_cpu_seconds <- df$seconds[which.max(cpu_col)] # Get the corresponding seconds
  
  # Append the max CPU value and the corresponding seconds to the data frame
  max_cpu_df <- rbind(max_cpu_df, data.frame(seconds = max_cpu_seconds, max_cpu = max_cpu))
}

max_cpu_df

# Add columns with average, minimum, and maximum values
df$minimum <- apply(df[, -1], 1, min) # Exclude seconds column from calculations
df$mean <- rowMeans(df[, -1])
df$maximum <- apply(df[, -1], 1, max)

# Rename the columns
colnames(df) <- c("seconds", "cpu-1", "cpu-2", "cpu-3", "cpu-4", "cpu-5", "minimum", "mean", "maximum")

# Plot
ggplot(df, aes(x = seconds)) +
  geom_line(aes(y = `cpu-1`, color = "CPU-1")) +
  geom_line(aes(y = `cpu-2`, color = "CPU-2")) +
  geom_line(aes(y = `cpu-3`, color = "CPU-3")) +
  geom_line(aes(y = `cpu-4`, color = "CPU-4")) +
  geom_line(aes(y = `cpu-5`, color = "CPU-5")) +
  xlab("Time Elapsed (s)") +
  ylab("CPU") +
  ggtitle(paste("CPU usage over time: ", test_name, sep="")) +
  theme(plot.title = element_text(hjust = 0.5)) +
  scale_color_manual(values = c("red", "green", "blue", "purple", "orange"), 
                     name = "Test Runs:", 
                     labels = c("1", "2", "3", "4", "5")) + 
  scale_x_continuous(breaks = seq(0, s, 20), limits = c(0, s)) +
  scale_y_continuous(breaks = seq(0, 0.01625, 0.0025), limits = c(0, 0.01625)) +
  theme_bw() +
  theme(plot.title = element_text(hjust = 0.5, size = 14), 
        plot.title.position = "plot", legend.position="bottom")




# Memory ----------------------------------------------------------------------
s <- 100

# Create an empty data frame with seconds 1 to 180
df <- data.frame(seconds = seq(1, s))

# Create an empty data frame to store max memory values
max_mem_df <- data.frame(seconds = numeric(), max_mem = numeric()) 

# Loop through the files and add the memory column to the data frame
for (i in 1:5) {
  file_name <- paste0(path, "mem-", i, ".csv")
  mem_data <- read.csv(file_name)
  mem_col <- mem_data[, 2][1:s] / 1048576 # Convert bytes to mebibytes
  df <- cbind(df, mem_col)

  max_mem <- max(mem_col) # Get the max memory value
  max_mem_seconds <- df$seconds[which.max(mem_col)] # Get the corresponding seconds
  
  # Append the max memory value and the corresponding seconds to the data frame
  max_mem_df <- rbind(max_mem_df, data.frame(seconds = max_mem_seconds, max_mem = max_mem))
}

max_mem_df

# Add columns with average, minimum, and maximum values
df$minimum <- apply(df[, -1], 1, min) # Exclude seconds column from calculations
df$mean <- rowMeans(df[, -1])
df$maximum <- apply(df[, -1], 1, max)

# Rename the columns
colnames(df) <- c("seconds", "mem-1", "mem-2", "mem-3", "mem-4", "mem-5", "minimum", "mean", "maximum")

# Plot
ggplot(df, aes(x = seconds)) +
  geom_line(aes(y = `mem-1`, color = "mem-1")) +
  geom_line(aes(y = `mem-2`, color = "mem-2")) +
  geom_line(aes(y = `mem-3`, color = "mem-3")) +
  geom_line(aes(y = `mem-4`, color = "mem-4")) +
  geom_line(aes(y = `mem-5`, color = "mem-5")) +
  xlab("Time Elapsed (s)") +
  ylab("Memory (MiB)") +
  ggtitle(paste("Memory usage over time: ", test_name, sep="")) +
  theme(plot.title = element_text(hjust = 0.5)) +
  scale_color_manual(values = c("red", "green", "blue", "purple", "orange"), 
                     name = "Test Runs:", 
                     labels = c("1", "2", "3", "4", "5")) + 
  scale_x_continuous(breaks = seq(0, s, 20), limits = c(0, s)) +
  scale_y_continuous(breaks = seq(0, 50, 10), limits = c(0, 50)) +
  theme_bw() +
  theme(plot.title = element_text(hjust = 0.5, size = 14), 
        plot.title.position = "plot", legend.position="bottom")









