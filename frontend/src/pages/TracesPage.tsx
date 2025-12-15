import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/Card'
import { Badge } from '../components/ui/Badge'
import { 
  Search,
  Filter,
  RefreshCw,
  ArrowUpDown,
  Eye,
  Clock,
  DollarSign,
  Zap
} from 'lucide-react'
import api from '../services/api'
import { formatCurrency, formatDuration, formatRelativeTime, getStatusColor } from '../lib/utils'

interface Trace {
  trace_id: string
  organization_id: string
  project_id: string
  timestamp: string
  trace_type: string
  duration_ms: number
  status: string
  total_cost_usd: number
  total_tokens: number
  model: string
  provider: string
  user_id?: string
}

export default function TracesPage() {
  const navigate = useNavigate()
  const [traces, setTraces] = useState<Trace[]>([])
  const [loading, setLoading] = useState(true)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [modelFilter, setModelFilter] = useState<string>('all')

  useEffect(() => {
    loadTraces()
  }, [statusFilter, modelFilter])

  const loadTraces = async () => {
    setLoading(true)
    try {
      const params: any = {
        organization_id: 'org-test',
        project_id: 'proj-test',
        limit: 50,
        offset: 0,
      }

      if (statusFilter !== 'all') {
        params.status = statusFilter
      }

      if (modelFilter !== 'all') {
        params.model = modelFilter
      }

      const response = await api.getTraces(params)
      console.log('Traces:', response)
      
      // Handle both possible response structures
      const tracesData = response.data?.traces || response.data || []
      setTraces(tracesData)
    } catch (error) {
      console.error('Failed to load traces:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredTraces = traces.filter(trace => {
    if (!searchTerm) return true
    const search = searchTerm.toLowerCase()
    return (
      trace.trace_id.toLowerCase().includes(search) ||
      trace.model.toLowerCase().includes(search) ||
      trace.status.toLowerCase().includes(search) ||
      (trace.user_id && trace.user_id.toLowerCase().includes(search))
    )
  })

  // Get unique models and statuses for filters
  const uniqueModels = Array.from(new Set(traces.map(t => t.model)))
  const uniqueStatuses = Array.from(new Set(traces.map(t => t.status)))

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Traces</h1>
          <p className="text-muted-foreground mt-2">
            View and debug your LLM traces
          </p>
        </div>

        <button
          onClick={loadTraces}
          className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          <RefreshCw className="h-4 w-4" />
          Refresh
        </button>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="grid gap-4 md:grid-cols-3">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search traces..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
              />
            </div>

            {/* Status Filter */}
            <div className="relative">
              <Filter className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
              >
                <option value="all">All Statuses</option>
                {uniqueStatuses.map(status => (
                  <option key={status} value={status}>{status}</option>
                ))}
              </select>
            </div>

            {/* Model Filter */}
            <div className="relative">
              <Filter className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <select
                value={modelFilter}
                onChange={(e) => setModelFilter(e.target.value)}
                className="w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-primary bg-background"
              >
                <option value="all">All Models</option>
                {uniqueModels.map(model => (
                  <option key={model} value={model}>{model}</option>
                ))}
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Stats Summary */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-2">
              <Zap className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total</span>
            </div>
            <div className="text-2xl font-bold mt-2">{filteredTraces.length}</div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-2">
              <DollarSign className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Cost</span>
            </div>
            <div className="text-2xl font-bold mt-2">
              {formatCurrency(filteredTraces.reduce((sum, t) => sum + t.total_cost_usd, 0))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Avg Latency</span>
            </div>
            <div className="text-2xl font-bold mt-2">
              {filteredTraces.length > 0 
                ? formatDuration(filteredTraces.reduce((sum, t) => sum + t.duration_ms, 0) / filteredTraces.length)
                : '0ms'
              }
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-2">
              <Zap className="h-4 w-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Success Rate</span>
            </div>
            <div className="text-2xl font-bold mt-2">
              {filteredTraces.length > 0
                ? ((filteredTraces.filter(t => t.status === 'success').length / filteredTraces.length) * 100).toFixed(1)
                : '0'
              }%
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Traces Table */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Traces ({filteredTraces.length})</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : filteredTraces.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              No traces found
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b text-left">
                    <th className="p-3 text-sm font-medium text-muted-foreground">
                      <div className="flex items-center gap-1">
                        Trace ID
                        <ArrowUpDown className="h-3 w-3" />
                      </div>
                    </th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Model</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Status</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Duration</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Cost</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Tokens</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Time</th>
                    <th className="p-3 text-sm font-medium text-muted-foreground">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredTraces.map((trace) => (
                    <tr 
                      key={trace.trace_id} 
                      className="border-b last:border-0 hover:bg-muted/50 cursor-pointer"
                      onClick={() => navigate(`/traces/${trace.trace_id}`)}
                    >
                      <td className="p-3">
                        <div className="font-mono text-sm">
                          {trace.trace_id.substring(0, 8)}...
                        </div>
                      </td>
                      <td className="p-3">
                        <div className="flex flex-col">
                          <span className="text-sm font-medium">{trace.model}</span>
                          <span className="text-xs text-muted-foreground">{trace.provider}</span>
                        </div>
                      </td>
                      <td className="p-3">
                        <Badge 
                          variant={trace.status === 'success' ? 'success' : 'destructive'}
                          className={getStatusColor(trace.status)}
                        >
                          {trace.status}
                        </Badge>
                      </td>
                      <td className="p-3 text-sm">{formatDuration(trace.duration_ms)}</td>
                      <td className="p-3 text-sm">{formatCurrency(trace.total_cost_usd)}</td>
                      <td className="p-3 text-sm">{trace.total_tokens.toLocaleString()}</td>
                      <td className="p-3 text-sm text-muted-foreground">
                        {formatRelativeTime(trace.timestamp)}
                      </td>
                      <td className="p-3">
                        <button
                          onClick={(e) => {
                            e.stopPropagation()
                            navigate(`/traces/${trace.trace_id}`)
                          }}
                          className="p-1.5 hover:bg-accent rounded-md"
                        >
                          <Eye className="h-4 w-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}