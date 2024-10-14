let userAddress = null;
let provider = null;
let signer = null;

window.onload = function() {
    const connectButton = document.getElementById("connectButton");
    const authButton = document.getElementById("authButton");
    const statusElement = document.getElementById("status");

    // Connect to MetaMask on button click
    connectButton.onclick = async function() {
        try {
            // Check if MetaMask is installed
            if (!window.ethereum) {
                alert("MetaMask not detected. Please install MetaMask!");
                return;
            }

            // Request account access
            provider = new ethers.BrowserProvider(window.ethereum);
            signer = await provider.getSigner();

            // Get and display the user's Ethereum address
            userAddress = await signer.getAddress();
            statusElement.textContent = `Connected: ${userAddress}`;

            // Enable the authentication button
            authButton.disabled = false;
        } catch (error) {
            console.error("Error connecting to MetaMask:", error);
        }
    };

    // Authenticate with the backend on button click
    authButton.onclick = async function() {
        try {
            // Request a nonce from the backend (replace with your own backend URL)
            const response = await fetch("http://localhost:8080/nonce", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ address: userAddress }),
            });
            const { nonce } = await response.json();

            // Sign the nonce using MetaMask
            const signature = await signer.signMessage(nonce);
            console.log("Signature:", signature);

            // Verify the signature with the backend
            const verifyResponse = await fetch("http://localhost:8080/auth", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ address: userAddress, signature }),
            });
            const { access_token } = await verifyResponse.json();

            // Display the token (or store it in localStorage for future use)
            console.log("Access Token:", access_token);
            statusElement.textContent = `Authenticated. Token: ${access_token}`;
        } catch (error) {
            console.error("Error during authentication:", error);
            statusElement.textContent = "Authentication failed.";
        }
    };
};