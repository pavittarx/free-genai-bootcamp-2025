FROM vllm/vllm-openai:latest

# Set environment variables
ENV HUGGING_FACE_HUB_TOKEN=""

# Create cache directory
RUN mkdir -p /root/.cache/huggingface

# Expose the port
EXPOSE 8000

# Set default command
CMD ["--model", "meta-llama/Llama-3.2-1B-Instruct", "--dtype", "half", "--gpu-memory-utilization", "0.8", "--swap-space", "2"]