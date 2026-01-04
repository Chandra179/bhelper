// src/components/money/money.js

export const moneyLogic = () => ({
    userInput: JSON.stringify({
        "example_1": "0.445454",
        "example_2": 0.342535,
        "example_3": 34343
    }, null, 2),
    results: [],

    processJson() {
        try {
            const data = JSON.parse(this.userInput);
            this.results = Object.entries(data).map(([key, value]) => {
                let units = 0;
                let detectedType = "";

                const isDecimal = value.toString().includes('.');

                if (isDecimal) {
                    detectedType = "Decimal (Converted)";
                    units = Math.round(parseFloat(value) * 100);
                } else {
                    detectedType = "Smallest Unit (Raw)";
                    units = parseInt(value);
                }

                return {
                    key: key,
                    units: units,
                    detected: detectedType,
                    formatted: (units / 100).toLocaleString('en-US', {
                        style: 'currency',
                        currency: 'USD'
                    })
                };
            });
        } catch (e) {
            alert("Invalid JSON");
        }
    }
});