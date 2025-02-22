import axios from "axios";
import {
  ConnectionInfo,
  InstanceId,
  INSTANCE_OPERATION_TIMEOUT,
  QueryInfo,
  ResourceObject,
  SqlResultSet,
} from "../../types";

function convert(resultSet: ResourceObject): SqlResultSet {
  return {
    data: JSON.parse((resultSet.attributes.data as string) || "{}"),
    error: resultSet.attributes.error as string,
  };
}

const getters = {
  convert:
    () =>
    (resultSet: ResourceObject): SqlResultSet => {
      return convert(resultSet);
    },
};

const actions = {
  async ping({ commit }: any, connectionInfo: ConnectionInfo) {
    const data = (
      await axios.post(`/api/sql/ping`, {
        data: {
          type: "connectionInfo",
          attributes: connectionInfo,
        },
      })
    ).data.data;

    return convert(data);
  },
  async syncSchema({ dispatch }: any, instanceId: InstanceId) {
    const data = (
      await axios.post(
        `/api/sql/syncschema`,
        {
          data: {
            type: "sqlSyncSchema",
            attributes: {
              instanceId,
            },
          },
        },
        {
          timeout: INSTANCE_OPERATION_TIMEOUT,
        }
      )
    ).data.data;

    const resultSet = convert(data);
    if (!resultSet.error) {
      // Refresh the corresponding list.
      dispatch("database/fetchDatabaseListByInstanceId", instanceId, {
        root: true,
      });

      dispatch("instance/fetchInstanceUserListById", instanceId, {
        root: true,
      });
    }

    return resultSet;
  },
  async query({ dispatch }: any, queryInfo: QueryInfo) {
    const data = (
      await axios.post(
        `/api/sql/execute`,
        {
          data: {
            type: "sqlExecute",
            attributes: {
              ...queryInfo,
              readonly: true,
            },
          },
        },
        {
          timeout: INSTANCE_OPERATION_TIMEOUT,
        }
      )
    ).data.data;

    const resultSet = convert(data);
    return resultSet.data;
  },
};

const mutations = {};

export default {
  namespaced: true,
  getters,
  actions,
  mutations,
};
