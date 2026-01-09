#!/bin/bash
# Load environment variables from .env file safely

# Function to load .env file
load_env() {
    if [ -f .env ]; then
        # Read and export all non-comment lines
        while IFS='=' read -r key value; do
            # Skip comments and empty lines
            [[ $key =~ ^[[:space:]]*# ]] && continue
            [[ -z $key ]] && continue
            
            # Remove leading/trailing whitespace
            key=$(echo "$key" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            value=$(echo "$value" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            
            # Remove quotes if present
            value=$(echo "$value" | sed 's/^"\(.*\)"$/\1/')
            
            # Export the variable
            export "$key=$value"
        done < .env
        
        echo "Environment variables loaded successfully"
    else
        echo "Error: .env file not found"
        exit 1
    fi
}

# Call the function
load_env
