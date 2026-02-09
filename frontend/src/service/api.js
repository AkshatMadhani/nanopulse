const API_BASE = 'http://localhost:8080'; 
class APIService {
 
  async fetch(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`;
    console.log('API Request:', url); 
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      });

      console.log('API Response status:', response.status); 
      if (!response.ok) {
        const error = await response.json().catch(() => ({ 
          message: response.statusText 
        }));
        throw new Error(error.message || `HTTP ${response.status}`);
      }

      const data = await response.json();
      console.log('API Response data:', data); 
      return data;
    } catch (error) {
      console.error(`API Error [${endpoint}]:`, error);
      throw error;
    }
  }

  async placeOrder(order) {
    return this.fetch('/order', {
      method: 'POST',
      body: JSON.stringify(order),
    });
  }

  async getHealth() {
    return this.fetch('/health');
  }

  async getStats() {
    return this.fetch('/stats');
  }

 
  async getOrderBook(symbol) {
    return this.fetch(`/book/${symbol}`);
  }
}

const apiService = new APIService();

export default apiService;