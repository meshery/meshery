/**
 * Very connection status
 * If connection already exist response will be connection object
 * If connection is not exist response will be 200 empty
 * 
 * @param {string} connectionKind - Kind of connection e.g (helm, github) pass only managed connections
*/
export const VerifyConnectionKind = async (connectionKind) => {
    return await axios.post(`/api/integrations/${connectionKind}/status`);
};