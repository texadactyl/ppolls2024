set -e

ppolls2024 -f
if [ $? -ne 0 ]; then
    exit 1
fi

ppolls2024 -l

ppolls2024 -p

ppolls2024 -r ec

