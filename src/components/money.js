// src/components/money.js

export const moneyTemplate = `
    <div class="max-w-6xl mx-auto p-4" x-data="moneyConverter">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="flex flex-col">
                <textarea 
                    x-model="userInput" 
                    class="w-full h-96 p-4 font-mono text-sm border-2 border-gray-300 rounded bg-gray-50 outline-none focus:border-blue-500"
                ></textarea>
                <button 
                    @click="processJson()"
                    class="mt-2 bg-black text-white font-bold py-2 rounded hover:bg-gray-800 transition"
                >
                    Process Mixed Units
                </button>
            </div>

            <div class="h-96 border-2 border-gray-300 rounded bg-white overflow-y-auto p-4 space-y-2">
                <template x-for="(item, index) in results" :key="index">
                    <div class="flex justify-between p-2 border-b font-mono text-sm items-center">
                        <div>
                            <span class="text-gray-500 block text-xs" x-text="item.key"></span>
                            <span class="text-[10px] text-gray-400" x-text="'Detected: ' + item.detected"></span>
                        </div>
                        <div class="text-right">
                            <span class="font-bold text-black" x-text="item.formatted"></span>
                            <span class="text-[10px] text-blue-600 block uppercase" x-text="item.units + ' units'"></span>
                        </div>
                    </div>
                </template>
            </div>
        </div>
    </div>
`;

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

                // Logic: If it contains a dot, it's a float/decimal (Dollar)
                // If no dot, it's an integer (Smallest Unit/Cents)
                const isDecimal = value.toString().includes('.');

                if (isDecimal) {
                    detectedType = "Decimal (Converted)";
                    // Convert dollar float to units (cents)
                    units = Math.round(parseFloat(value) * 100);
                } else {
                    detectedType = "Smallest Unit (Raw)";
                    // Treat as raw units
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