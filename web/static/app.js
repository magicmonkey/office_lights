// State management
let currentState = null;
let updateTimer = null;
let pollTimer = null;
let isUpdating = false;

// Debounce delay in milliseconds
const DEBOUNCE_DELAY = 300;
const POLL_INTERVAL = 3000;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initializeEventListeners();
    loadInitialState();
    startPolling();
});

// Initialize all event listeners
function initializeEventListeners() {
    // LED Strip
    document.getElementById('strip-r').addEventListener('input', handleStripRGBChange);
    document.getElementById('strip-g').addEventListener('input', handleStripRGBChange);
    document.getElementById('strip-b').addEventListener('input', handleStripRGBChange);
    document.getElementById('strip-color').addEventListener('input', handleStripColorPickerChange);

    // LED Bar - Section buttons
    document.getElementById('ledbar-section-1').addEventListener('click', () => handleLEDBarSection(1));
    document.getElementById('ledbar-section-2').addEventListener('click', () => handleLEDBarSection(2));

    // LED Bar - Mode buttons
    document.getElementById('ledbar-mode-rgbw').addEventListener('click', () => handleLEDBarMode('rgbw'));
    document.getElementById('ledbar-mode-white').addEventListener('click', () => handleLEDBarMode('white'));

    // LED Bar - RGBW controls
    document.getElementById('ledbar-led').addEventListener('input', handleLEDBarLEDChange);
    document.getElementById('ledbar-r').addEventListener('input', handleLEDBarRGBWChange);
    document.getElementById('ledbar-g').addEventListener('input', handleLEDBarRGBWChange);
    document.getElementById('ledbar-b').addEventListener('input', handleLEDBarRGBWChange);
    document.getElementById('ledbar-w').addEventListener('input', handleLEDBarRGBWChange);

    // LED Bar - White controls
    document.getElementById('ledbar-white-led').addEventListener('input', handleLEDBarWhiteLEDChange);
    document.getElementById('ledbar-white').addEventListener('input', handleLEDBarWhiteChange);

    // Video Light 1
    document.getElementById('vl1-on').addEventListener('change', handleVideoLight1Change);
    document.getElementById('vl1-brightness').addEventListener('input', handleVideoLight1Change);

    // Video Light 2
    document.getElementById('vl2-on').addEventListener('change', handleVideoLight2Change);
    document.getElementById('vl2-brightness').addEventListener('input', handleVideoLight2Change);
}

// Load initial state from server
async function loadInitialState() {
    try {
        const response = await fetch('/api');
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        currentState = await response.json();
        updateUIFromState(currentState);
        updateConnectionStatus(true);
        updateLastUpdateTime();
        hideError();
    } catch (error) {
        console.error('Failed to load initial state:', error);
        showError('Failed to connect to server: ' + error.message);
        updateConnectionStatus(false);
    }
}

// Start polling for state updates
function startPolling() {
    pollTimer = setInterval(async () => {
        if (!isUpdating) {
            await loadInitialState();
        }
    }, POLL_INTERVAL);
}

// Update UI controls from state
function updateUIFromState(state) {
    if (!state) return;

    // LED Strip
    updateStripUI(state.ledStrip);

    // LED Bar - we need to update based on current section and mode
    const currentSection = getCurrentLEDBarSection();
    const currentMode = getCurrentLEDBarMode();
    updateLEDBarUI(state.ledBar, currentSection, currentMode);

    // Video Lights
    updateVideoLightUI(1, state.videoLight1);
    updateVideoLightUI(2, state.videoLight2);
}

// Update LED Strip UI
function updateStripUI(ledStrip) {
    document.getElementById('strip-r').value = ledStrip.r;
    document.getElementById('strip-g').value = ledStrip.g;
    document.getElementById('strip-b').value = ledStrip.b;

    document.getElementById('strip-r-value').textContent = ledStrip.r;
    document.getElementById('strip-g-value').textContent = ledStrip.g;
    document.getElementById('strip-b-value').textContent = ledStrip.b;

    // Update color picker
    const hexColor = rgbToHex(ledStrip.r, ledStrip.g, ledStrip.b);
    document.getElementById('strip-color').value = hexColor;

    // Update preview
    const preview = document.getElementById('strip-preview');
    preview.style.backgroundColor = `rgb(${ledStrip.r}, ${ledStrip.g}, ${ledStrip.b})`;
}

// Update LED Bar UI
function updateLEDBarUI(ledBar, section, mode) {
    const sectionData = section === 1 ? ledBar.section1 : ledBar.section2;

    if (mode === 'rgbw') {
        const ledIndex = parseInt(document.getElementById('ledbar-led').value) - 1;
        if (ledIndex >= 0 && ledIndex < sectionData.rgbw.length) {
            const rgbw = sectionData.rgbw[ledIndex];
            document.getElementById('ledbar-r').value = rgbw.r;
            document.getElementById('ledbar-g').value = rgbw.g;
            document.getElementById('ledbar-b').value = rgbw.b;
            document.getElementById('ledbar-w').value = rgbw.w;

            document.getElementById('ledbar-r-value').textContent = rgbw.r;
            document.getElementById('ledbar-g-value').textContent = rgbw.g;
            document.getElementById('ledbar-b-value').textContent = rgbw.b;
            document.getElementById('ledbar-w-value').textContent = rgbw.w;
        }
    } else {
        const ledIndex = parseInt(document.getElementById('ledbar-white-led').value) - 1;
        if (ledIndex >= 0 && ledIndex < sectionData.white.length) {
            const white = sectionData.white[ledIndex];
            document.getElementById('ledbar-white').value = white;
            document.getElementById('ledbar-white-value').textContent = white;
        }
    }
}

// Update Video Light UI
function updateVideoLightUI(lightNum, lightState) {
    const prefix = `vl${lightNum}`;

    document.getElementById(`${prefix}-on`).checked = lightState.on;
    document.getElementById(`${prefix}-brightness`).value = lightState.brightness;
    document.getElementById(`${prefix}-brightness-value`).textContent = lightState.brightness;

    // Update indicator
    const indicator = document.getElementById(`${prefix}-indicator`);
    if (lightState.on) {
        indicator.classList.add('on');
        const brightness = lightState.brightness / 100;
        indicator.style.opacity = 0.3 + (brightness * 0.7);
    } else {
        indicator.classList.remove('on');
        indicator.style.opacity = 1;
    }
}

// LED Strip RGB slider change
function handleStripRGBChange() {
    const r = parseInt(document.getElementById('strip-r').value);
    const g = parseInt(document.getElementById('strip-g').value);
    const b = parseInt(document.getElementById('strip-b').value);

    document.getElementById('strip-r-value').textContent = r;
    document.getElementById('strip-g-value').textContent = g;
    document.getElementById('strip-b-value').textContent = b;

    // Update color picker
    const hexColor = rgbToHex(r, g, b);
    document.getElementById('strip-color').value = hexColor;

    // Update preview
    const preview = document.getElementById('strip-preview');
    preview.style.backgroundColor = `rgb(${r}, ${g}, ${b})`;

    // Update state
    currentState.ledStrip.r = r;
    currentState.ledStrip.g = g;
    currentState.ledStrip.b = b;

    debouncedUpdate();
}

// LED Strip color picker change
function handleStripColorPickerChange() {
    const hexColor = document.getElementById('strip-color').value;
    const rgb = hexToRgb(hexColor);

    document.getElementById('strip-r').value = rgb.r;
    document.getElementById('strip-g').value = rgb.g;
    document.getElementById('strip-b').value = rgb.b;

    document.getElementById('strip-r-value').textContent = rgb.r;
    document.getElementById('strip-g-value').textContent = rgb.g;
    document.getElementById('strip-b-value').textContent = rgb.b;

    // Update preview
    const preview = document.getElementById('strip-preview');
    preview.style.backgroundColor = hexColor;

    // Update state
    currentState.ledStrip.r = rgb.r;
    currentState.ledStrip.g = rgb.g;
    currentState.ledStrip.b = rgb.b;

    debouncedUpdate();
}

// LED Bar section button
function handleLEDBarSection(section) {
    document.getElementById('ledbar-section-1').classList.toggle('active', section === 1);
    document.getElementById('ledbar-section-2').classList.toggle('active', section === 2);

    // Update UI with new section data
    const mode = getCurrentLEDBarMode();
    updateLEDBarUI(currentState.ledBar, section, mode);
}

// LED Bar mode button
function handleLEDBarMode(mode) {
    document.getElementById('ledbar-mode-rgbw').classList.toggle('active', mode === 'rgbw');
    document.getElementById('ledbar-mode-white').classList.toggle('active', mode === 'white');

    // Show/hide controls
    document.getElementById('ledbar-rgbw-controls').style.display = mode === 'rgbw' ? 'block' : 'none';
    document.getElementById('ledbar-white-controls').style.display = mode === 'white' ? 'block' : 'none';

    // Update UI with current mode data
    const section = getCurrentLEDBarSection();
    updateLEDBarUI(currentState.ledBar, section, mode);
}

// LED Bar LED selector change (RGBW)
function handleLEDBarLEDChange() {
    const section = getCurrentLEDBarSection();
    updateLEDBarUI(currentState.ledBar, section, 'rgbw');
}

// LED Bar RGBW slider change
function handleLEDBarRGBWChange() {
    const r = parseInt(document.getElementById('ledbar-r').value);
    const g = parseInt(document.getElementById('ledbar-g').value);
    const b = parseInt(document.getElementById('ledbar-b').value);
    const w = parseInt(document.getElementById('ledbar-w').value);

    document.getElementById('ledbar-r-value').textContent = r;
    document.getElementById('ledbar-g-value').textContent = g;
    document.getElementById('ledbar-b-value').textContent = b;
    document.getElementById('ledbar-w-value').textContent = w;

    // Update state
    const section = getCurrentLEDBarSection();
    const ledIndex = parseInt(document.getElementById('ledbar-led').value) - 1;
    const sectionData = section === 1 ? currentState.ledBar.section1 : currentState.ledBar.section2;

    if (ledIndex >= 0 && ledIndex < sectionData.rgbw.length) {
        sectionData.rgbw[ledIndex] = { r, g, b, w };
        debouncedUpdate();
    }
}

// LED Bar white LED selector change
function handleLEDBarWhiteLEDChange() {
    const section = getCurrentLEDBarSection();
    updateLEDBarUI(currentState.ledBar, section, 'white');
}

// LED Bar white slider change
function handleLEDBarWhiteChange() {
    const white = parseInt(document.getElementById('ledbar-white').value);
    document.getElementById('ledbar-white-value').textContent = white;

    // Update state
    const section = getCurrentLEDBarSection();
    const ledIndex = parseInt(document.getElementById('ledbar-white-led').value) - 1;
    const sectionData = section === 1 ? currentState.ledBar.section1 : currentState.ledBar.section2;

    if (ledIndex >= 0 && ledIndex < sectionData.white.length) {
        sectionData.white[ledIndex] = white;
        debouncedUpdate();
    }
}

// Video Light 1 change
function handleVideoLight1Change() {
    const on = document.getElementById('vl1-on').checked;
    const brightness = parseInt(document.getElementById('vl1-brightness').value);

    document.getElementById('vl1-brightness-value').textContent = brightness;

    // Update indicator
    const indicator = document.getElementById('vl1-indicator');
    if (on) {
        indicator.classList.add('on');
        const brightnessValue = brightness / 100;
        indicator.style.opacity = 0.3 + (brightnessValue * 0.7);
    } else {
        indicator.classList.remove('on');
        indicator.style.opacity = 1;
    }

    // Update state
    currentState.videoLight1.on = on;
    currentState.videoLight1.brightness = brightness;

    debouncedUpdate();
}

// Video Light 2 change
function handleVideoLight2Change() {
    const on = document.getElementById('vl2-on').checked;
    const brightness = parseInt(document.getElementById('vl2-brightness').value);

    document.getElementById('vl2-brightness-value').textContent = brightness;

    // Update indicator
    const indicator = document.getElementById('vl2-indicator');
    if (on) {
        indicator.classList.add('on');
        const brightnessValue = brightness / 100;
        indicator.style.opacity = 0.3 + (brightnessValue * 0.7);
    } else {
        indicator.classList.remove('on');
        indicator.style.opacity = 1;
    }

    // Update state
    currentState.videoLight2.on = on;
    currentState.videoLight2.brightness = brightness;

    debouncedUpdate();
}

// Debounced update to server
function debouncedUpdate() {
    if (updateTimer) {
        clearTimeout(updateTimer);
    }

    updateTimer = setTimeout(() => {
        sendStateToServer();
    }, DEBOUNCE_DELAY);
}

// Send state to server
async function sendStateToServer() {
    if (!currentState) return;

    isUpdating = true;

    try {
        const response = await fetch('/api', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(currentState),
        });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const updatedState = await response.json();
        currentState = updatedState;
        updateConnectionStatus(true);
        updateLastUpdateTime();
        hideError();
    } catch (error) {
        console.error('Failed to update state:', error);
        showError('Failed to update lights: ' + error.message);
        updateConnectionStatus(false);
    } finally {
        isUpdating = false;
    }
}

// Get current LED Bar section
function getCurrentLEDBarSection() {
    return document.getElementById('ledbar-section-1').classList.contains('active') ? 1 : 2;
}

// Get current LED Bar mode
function getCurrentLEDBarMode() {
    return document.getElementById('ledbar-mode-rgbw').classList.contains('active') ? 'rgbw' : 'white';
}

// Update connection status
function updateConnectionStatus(connected) {
    const status = document.getElementById('connection-status');
    if (connected) {
        status.textContent = 'Connected';
        status.classList.remove('disconnected');
        status.classList.add('connected');
    } else {
        status.textContent = 'Disconnected';
        status.classList.remove('connected');
        status.classList.add('disconnected');
    }
}

// Update last update time
function updateLastUpdateTime() {
    const now = new Date();
    const timeString = now.toLocaleTimeString();
    document.getElementById('last-update').textContent = `Last update: ${timeString}`;
}

// Show error message
function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}

// Hide error message
function hideError() {
    const errorDiv = document.getElementById('error-message');
    errorDiv.style.display = 'none';
}

// Convert RGB to hex
function rgbToHex(r, g, b) {
    return '#' + [r, g, b].map(x => {
        const hex = x.toString(16);
        return hex.length === 1 ? '0' + hex : hex;
    }).join('');
}

// Convert hex to RGB
function hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : { r: 0, g: 0, b: 0 };
}
