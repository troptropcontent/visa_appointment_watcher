# Install KAMAL
echo "Installing KAMAL..."
gem install kamal
echo "KAMAL installed."

# Install tailwindcss if it is not there yet
echo "Installing tailwindcss..."
FILE=./bin/tailwindcss
if [ -f "$FILE" ]; then
    echo "$FILE already exists, does not need to be installed."
else 
    echo "$FILE does not exist"
    echo "Downloading tailwindcss..."
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
    echo "Installing tailwindcss..."
    chmod +x tailwindcss-linux-arm64
    sudo mv tailwindcss-linux-arm64 ./bin/tailwindcss
    echo "tailwindcss installed."
fi

