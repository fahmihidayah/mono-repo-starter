'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { useState } from 'react'
import { toast } from 'sonner'
import { logger } from '~/utils/logger'

const queryLogger = logger.scope('ReactQuery')

// Advanced QueryClient configuration
function createQueryClient() {
  queryLogger.info('Creating QueryClient...')

  return new QueryClient({
    defaultOptions: {
      queries: {
        // Stale time: How long until data is considered stale
        staleTime: 60 * 1000, // 1 minute

        // Garbage collection time: How long unused data stays in cache
        gcTime: 5 * 60 * 1000, // 5 minutes

        // Retry failed requests
        retry: (failureCount, error: any) => {
          queryLogger.debug('Query retry check', {
            failureCount,
            errorStatus: error?.status,
          })

          // Don't retry on 4xx errors (client errors)
          if (error?.status >= 400 && error?.status < 500) {
            queryLogger.warn('Not retrying 4xx error', { status: error.status })
            return false
          }
          // Retry up to 3 times for 5xx errors
          const shouldRetry = failureCount < 3
          queryLogger.info(`Query retry: ${shouldRetry}`, { failureCount })
          return shouldRetry
        },

        // Retry delay with exponential backoff
        retryDelay: (attemptIndex) => {
          const delay = Math.min(1000 * 2 ** attemptIndex, 30000)
          queryLogger.debug('Query retry delay', { attemptIndex, delay })
          return delay
        },

        // Refetch on window focus for fresh data
        refetchOnWindowFocus: true,

        // Don't refetch on mount if data is fresh
        refetchOnMount: true,

        // Refetch on reconnect
        refetchOnReconnect: true,

        // Network mode
        networkMode: 'online',
      },
      mutations: {
        // Network mode for mutations
        networkMode: 'online',

        // Retry mutations once on network error
        retry: (failureCount, error: any) => {
          queryLogger.debug('Mutation retry check', {
            failureCount,
            errorStatus: error?.status,
          })

          // Don't retry on 4xx errors
          if (error?.status >= 400 && error?.status < 500) {
            queryLogger.warn('Not retrying 4xx mutation error', { status: error.status })
            return false
          }
          // Retry once for network errors
          const shouldRetry = failureCount < 1
          queryLogger.info(`Mutation retry: ${shouldRetry}`, { failureCount })
          return shouldRetry
        },

        // Global error handler for mutations
        onError: (error: any, _variables, _context) => {
          const errorMessage =
            error?.message ||
            error?.response?.data?.message ||
            'An error occurred. Please try again.'

          queryLogger.error('Mutation error', {
            message: errorMessage,
            error,
          })

          toast.error('Error', {
            description: errorMessage,
          })
        },

        // Global success handler for mutations
        onSuccess: (_data, _variables, _context) => {
          queryLogger.info('Mutation success')
          // You can enable this for global success toasts
          // toast.success('Success', {
          //   description: 'Operation completed successfully',
          // })
        },
      },
    },
  })
}

export function QueryProvider({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => createQueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {/* React Query DevTools - Only visible in development */}
      <ReactQueryDevtools
        initialIsOpen={false}
        position="bottom"
        buttonPosition="bottom-right"
      />
    </QueryClientProvider>
  )
}