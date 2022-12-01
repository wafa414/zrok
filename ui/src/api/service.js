/** @module service */
// Auto-generated, edits will be overwritten
import * as gateway from './gateway'

/**
 * @param {object} options Optional options
 * @param {module:types.accessRequest} [options.body] 
 * @return {Promise<module:types.accessResponse>} access created
 */
export function access(options) {
  if (!options) options = {}
  const parameters = {
    body: {
      body: options.body
    }
  }
  return gateway.request(accessOperation, parameters)
}

/**
 * @param {object} options Optional options
 * @param {module:types.serviceRequest} [options.body] 
 * @return {Promise<module:types.service>} ok
 */
export function getService(options) {
  if (!options) options = {}
  const parameters = {
    body: {
      body: options.body
    }
  }
  return gateway.request(getServiceOperation, parameters)
}

/**
 * @param {object} options Optional options
 * @param {module:types.shareRequest} [options.body] 
 * @return {Promise<module:types.shareResponse>} service created
 */
export function share(options) {
  if (!options) options = {}
  const parameters = {
    body: {
      body: options.body
    }
  }
  return gateway.request(shareOperation, parameters)
}

/**
 * @param {object} options Optional options
 * @param {module:types.unaccessRequest} [options.body] 
 * @return {Promise<object>} access removed
 */
export function unaccess(options) {
  if (!options) options = {}
  const parameters = {
    body: {
      body: options.body
    }
  }
  return gateway.request(unaccessOperation, parameters)
}

/**
 * @param {object} options Optional options
 * @param {module:types.unshareRequest} [options.body] 
 * @return {Promise<object>} service removed
 */
export function unshare(options) {
  if (!options) options = {}
  const parameters = {
    body: {
      body: options.body
    }
  }
  return gateway.request(unshareOperation, parameters)
}

const accessOperation = {
  path: '/access',
  contentTypes: ['application/zrok.v1+json'],
  method: 'post',
  security: [
    {
      id: 'key'
    }
  ]
}

const getServiceOperation = {
  path: '/service',
  contentTypes: ['application/zrok.v1+json'],
  method: 'get',
  security: [
    {
      id: 'key'
    }
  ]
}

const shareOperation = {
  path: '/share',
  contentTypes: ['application/zrok.v1+json'],
  method: 'post',
  security: [
    {
      id: 'key'
    }
  ]
}

const unaccessOperation = {
  path: '/unaccess',
  contentTypes: ['application/zrok.v1+json'],
  method: 'delete',
  security: [
    {
      id: 'key'
    }
  ]
}

const unshareOperation = {
  path: '/unshare',
  contentTypes: ['application/zrok.v1+json'],
  method: 'delete',
  security: [
    {
      id: 'key'
    }
  ]
}
