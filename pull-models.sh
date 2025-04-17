#!/bin/sh

echo "🟡 Waiting for Ollama to be ready..."
sleep 15  # Initial wait to let Ollama initialize

# More robust wait mechanism with timeout
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if curl -s -f http://ollama:11434/api/version > /dev/null 2>&1; then
    echo "🟢 Ollama is ready! API connection successful."
    break
  else
    RETRY_COUNT=$((RETRY_COUNT+1))
    echo "Attempt $RETRY_COUNT/$MAX_RETRIES: Ollama not ready yet, retrying in 10 seconds..."
    sleep 10
  fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
  echo "❌ Failed to connect to Ollama after $MAX_RETRIES attempts. Exiting."
  exit 1
fi

echo "🔄 Pulling models in order of size (smallest first)..."

# First tier - Smallest models (best for CPU)
small_models="
phi4-mini
qwen2.5:0.5b
qwen2.5-coder:0.5b
gemma2:2b
tinyllama
"

echo "📦 TIER 1: Pulling smallest models first (best for CPU operation)..."
for model in $small_models; do
  echo "👉 Pulling $model"
  response=$(curl -s -X POST http://ollama:11434/api/pull -d "{\"name\":\"$model\"}")
  echo "$response"
  
  if echo "$response" | grep -q "error"; then
    echo "⚠️ Warning: Issue with pulling $model - will try alternate version"
    # Try without specific tag if failed
    if [[ $model == *":"* ]]; then
      base_model=${model%%:*}
      echo "👉 Trying $base_model instead"
      response=$(curl -s -X POST http://ollama:11434/api/pull -d "{\"name\":\"$base_model\"}")
      echo "$response"
    fi
  else
    echo "✅ Model $model pull initiated successfully"
  fi
  
  sleep 5
done

# Second tier - Medium models
medium_models="
phi3
phi3.5
dolphin-phi
stable-code:3b
gemma2:2b
qwen2.5:1.5b
qwen2.5-coder:1.5b
qwen2-math:1.5b
orca-mini:3b-q4_0
codegemma:2b
"

echo "📦 TIER 2: Pulling medium-sized models..."
for model in $medium_models; do
  echo "👉 Pulling $model"
  response=$(curl -s -X POST http://ollama:11434/api/pull -d "{\"name\":\"$model\"}")
  echo "$response"
  
  if echo "$response" | grep -q "error"; then
    echo "⚠️ Warning: Issue with pulling $model - skipping"
  else
    echo "✅ Model $model pull initiated successfully"
  fi
  
  sleep 5
done

# Third tier - Larger models (may be very slow on CPU)
large_models="
neural-chat:7b-v3.3-q4_0
mistral:7b-instruct-v0.2-q4_0
wizard-math:7b-q4_0
openchat:7b-v3.5-q4_0
"

echo "📦 TIER 3: Pulling larger models (these will be slow on CPU)..."
for model in $large_models; do
  echo "👉 Pulling $model"
  response=$(curl -s -X POST http://ollama:11434/api/pull -d "{\"name\":\"$model\"}")
  echo "$response"
  
  if echo "$response" | grep -q "error"; then
    echo "⚠️ Warning: Issue with pulling $model - skipping"
  else
    echo "✅ Model $model pull initiated successfully"
  fi
  
  sleep 5
done

echo "✅ All models have been queued for download"
echo "NOTE: Models will continue downloading in the background."
echo "NOTE: On CPU-only systems, models larger than 3B parameters may be very slow."