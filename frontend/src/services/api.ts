import axios, { AxiosInstance } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const API_KEY = 'demo-key-456' // Use test API key for demo

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': API_KEY, // Use API key instead of JWT
      },
    })
  }

  async login(email: string, _password: string) {
    // Mock login - doesn't call backend
    return {
      data: {
        token: 'mock-token',
        user: {
          id: 'user-test-123',
          email: email,
          name: 'Diksha Sahare',
          organization_id: 'org-test-123',
          role: 'admin',
          created_at: new Date().toISOString(),
        }
      }
    }
  }

  async getCurrentUser() {
    return {
      data: {
        id: 'user-test-123',
        email: 'diksha@example.com',
        name: 'Diksha Sahare',
        organization_id: 'org-test-123',
        role: 'admin',
        created_at: new Date().toISOString(),
      }
    }
  }

  async getDashboard(timeRange: string = '24h') {
    const response = await this.client.get('/api/v1/analytics/dashboard', {
      params: { time_range: timeRange },
    })
    return response.data
  }

  async getTraces(params: any) {
    const response = await this.client.get('/api/v1/traces', { params })
    return response.data
  }

  async getTrace(traceId: string) {
    const response = await this.client.get(`/api/v1/traces/${traceId}`)
    return response.data
  }

  async checkHealth() {
    const response = await this.client.get('/health')
    return response.data
  }
}

export const api = new ApiClient()
export default api