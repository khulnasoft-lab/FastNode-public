import * as actions from '../actions/system'

const defaultState = {
  version: "",
  networkConnected: false,
  haveCheckedNetworkConnection: false,
  fastnodedReady: false,
  hasSeenSpyderNotification: false,
}

const update = (state, action) => ({
  ...state,
  ...action.data,
})

const updateVersion = (state, action) => ({
  ...state,
  version: action.data,
})

const updateNetworkConnected = (state, action) => ({
  ...state,
  networkConnected: action.isOnline,
  haveCheckedNetworkConnection: true,
})

const updateFastnodedReady = (state, action) => ({
  ...state,
  fastnodedReady: action.fastnodedReady,
})

const updateHasSeenSpyderNotification = (state, action) => ({
  ...state,
  hasSeenSpyderNotification: action.data === true,
})

const system = (state = defaultState, action) => {
  switch(action.type) {
    case actions.UPDATE_SYSTEM_INFO:
      return update(state, action)
    case actions.GET_VERSION:
      return updateVersion(state, action)
    case actions.CHECK_IF_ONLINE:
      return updateNetworkConnected(state, action)
    case actions.GET_FASTNODED_READY:
      return updateFastnodedReady(state, action)
    case actions.SET_HAS_SEEN_SPYDER_NOTIFICATION:
      return updateHasSeenSpyderNotification(state, action)
    default:
      return state
  }
}

export default system
