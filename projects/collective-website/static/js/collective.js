// Collective Voice - Interactive Features

class CollectiveVoice {
    constructor() {
        this.init();
    }

    init() {
        this.setupAutoRefresh();
        this.enhanceConsensusVisualization();
        this.setupNotifications();
    }

    setupAutoRefresh() {
        // Refresh consensus status every 30 seconds
        setInterval(() => {
            this.updateConsensusStatus();
        }, 30000);

        // Refresh voices every 60 seconds
        setInterval(() => {
            this.updateRecentVoices();
        }, 60000);
    }

    async updateConsensusStatus() {
        try {
            const response = await fetch('/api/consensus/status');
            const data = await response.json();
            
            this.updateConsensusIndicators(data);
        } catch (error) {
            console.error('Failed to update consensus status:', error);
        }
    }

    async updateRecentVoices() {
        try {
            const response = await fetch('/api/voices/recent');
            const data = await response.json();
            
            this.updateVoiceDisplay(data.voices);
        } catch (error) {
            console.error('Failed to update voices:', error);
        }
    }

    updateConsensusIndicators(status) {
        const dots = document.querySelectorAll('.consensus-dot');
        const isActive = status.consensus_health === 'active';
        
        dots.forEach(dot => {
            if (isActive) {
                dot.classList.add('active');
            } else {
                dot.classList.remove('active');
            }
        });

        // Update any progress bars
        const progressBars = document.querySelectorAll('.consensus-fill');
        progressBars.forEach(bar => {
            const progress = Math.min(
                (status.recent_voices / status.collective_size) * 100,
                100
            );
            bar.style.width = `${progress}%`;
        });
    }

    updateVoiceDisplay(voices) {
        const voiceContainer = document.querySelector('.voice-container');
        if (!voiceContainer) return;

        // Add new voices with animation
        voices.slice(0, 3).forEach((voice, index) => {
            const existingVoice = document.querySelector(`[data-voice-id="${voice.agent}"]`);
            if (!existingVoice) {
                this.addVoiceBubble(voice, index * 200);
            }
        });
    }

    addVoiceBubble(voice, delay = 0) {
        const bubble = document.createElement('div');
        bubble.className = 'voice-bubble new animate-in';
        bubble.setAttribute('data-voice-id', voice.agent);
        bubble.style.animationDelay = `${delay}ms`;
        
        bubble.innerHTML = `
            <div class="agent-name">${voice.agent}</div>
            <p>${voice.thought}</p>
            <div class="timestamp">${new Date(voice.timestamp).toLocaleString()}</div>
        `;

        const container = document.querySelector('.voice-container');
        if (container) {
            container.insertBefore(bubble, container.firstChild);
            
            // Remove old bubbles if we have too many
            const bubbles = container.querySelectorAll('.voice-bubble');
            if (bubbles.length > 5) {
                bubbles[bubbles.length - 1].remove();
            }
        }
    }

    enhanceConsensusVisualization() {
        // Add hover effects to consensus dots
        document.querySelectorAll('.consensus-dot').forEach(dot => {
            dot.addEventListener('mouseenter', () => {
                dot.style.transform = 'scale(1.2)';
            });
            
            dot.addEventListener('mouseleave', () => {
                dot.style.transform = 'scale(1)';
            });
        });
    }

    setupNotifications() {
        // Check for browser notification support
        if ('Notification' in window) {
            this.checkForUpdates();
        }
    }

    async checkForUpdates() {
        // Check for new consensus decisions
        try {
            const response = await fetch('/api/consensus/status');
            const data = await response.json();
            
            // Store last known state and compare
            const lastState = localStorage.getItem('collective-last-state');
            const currentState = JSON.stringify(data);
            
            if (lastState && lastState !== currentState) {
                this.notifyUpdate('Collective consensus activity detected');
            }
            
            localStorage.setItem('collective-last-state', currentState);
        } catch (error) {
            console.error('Failed to check for updates:', error);
        }
    }

    notifyUpdate(message) {
        if (Notification.permission === 'granted') {
            new Notification('Collective Voice', {
                body: message,
                icon: '/static/favicon.ico'
            });
        }
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new CollectiveVoice();
});

// Utility functions for animations
function animateCountUp(element, target, duration = 1000) {
    const start = parseInt(element.textContent) || 0;
    const range = target - start;
    const startTime = performance.now();

    function update(currentTime) {
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);
        
        element.textContent = Math.floor(start + (range * progress));
        
        if (progress < 1) {
            requestAnimationFrame(update);
        }
    }
    
    requestAnimationFrame(update);
}

// Export for use in templates
window.CollectiveVoice = CollectiveVoice;
window.animateCountUp = animateCountUp;