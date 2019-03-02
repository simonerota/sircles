import React from 'react'
import { graphql, compose } from 'react-apollo'
import gql from 'graphql-tag'
import { Container, Dimmer, Loader, Segment, Input, Button, Table, Icon, Confirm, Popup } from 'semantic-ui-react'

import { withError } from '../modules/Error'
import Util from '../modules/Util'
import Avatar from '../components/Avatar'

const defaultFetchSize = 25
let uuuid

class Members extends React.Component {
  componentWillMount () {
    this.resetComponent()
  }

  resetComponent = () => this.setState({
    searchString: '',
    isOpenActivateMember: false
  })

  componentWillReceiveProps (nextProps) {
    const { viewerQuery, membersQuery } = this.props

    if (Util.isQueriesError(viewerQuery, membersQuery)) {
      this.props.appError.setError(true)
      return
    }
  }

  handleSearchStringChange = (e, data) => {
    this.setState({searchString: data.value})
    console.log(this.props)
    this.props.membersQuery.update(data.value).catch((error) => {
      this.props.appError.setError(true)
      console.log('there was an error sending the query', error)
    })
  }

  handleNewMember = () => {
    this.props.history.push('/settings/admin/member/new')
  }

  handleActivateMember = () => {
      let updateActivateMemberChange =
        {
          uid: uuuid
        }

        console.log('updateActivateMemberChange', updateActivateMemberChange)
      this.props.updateActivateMember(updateActivateMemberChange)
    .then(({ data }) => {
      console.log("sending query ActivateMember")
    }).catch((error) => {
      console.log('there was an error sending the query', error)
    })
    location.reload()
    this.closeActivateMember()
  }

  openActivateMember = (uid) => {
    uuuid = uid
    this.setState({ 'isOpenActivateMember': true })
  }

  closeActivateMember = () => {
    this.setState({ 'isOpenActivateMember': false })
  }

  render () {
    const { viewerQuery, membersQuery } = this.props
    let loading = false

    console.log(this.props)

    if (Util.isQueriesError(viewerQuery, membersQuery)) {
      return null
    }

    if (Util.isQueriesLoading(viewerQuery)) {
      return (
        <Dimmer active inverted>
          <Loader inverted>Loading</Loader>
        </Dimmer>
      )
    }

    if (Util.isQueriesLoading(membersQuery)) {
      loading = true
    }

    console.log('membersQuery', membersQuery)

    const { searchString } = this.state

    const viewer = viewerQuery.viewer
    const members = membersQuery.membersList

    return (
      <Container>
        <Segment clearing>
          { viewer.member.isAdmin &&
          <Button floated='right' color='green' onClick={this.handleNewMember}>New Member</Button>
              }
          <div>
            <h2>Members</h2>
          </div>
        </Segment>

        { loading &&
        <Dimmer active inverted>
          <Loader inverted>Loading</Loader>
        </Dimmer>
        }

        { !loading &&
          <div>
            <Input icon='users' iconPosition='left' value={searchString} onChange={this.handleSearchStringChange} placeholder='Search users...' />
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>UserName</Table.HeaderCell>
                  <Table.HeaderCell>Full Name</Table.HeaderCell>
                  <Table.HeaderCell>Admin</Table.HeaderCell>
                  <Table.HeaderCell>Disabled</Table.HeaderCell>
                  <Table.HeaderCell />
                </Table.Row>
              </Table.Header>

              <Table.Body>
                {members.map(m => (
                  <Table.Row key={m.uid}>
                    <Table.Cell>
                      <Avatar uid={m.uid} size={40} inline spaced shape='rounded' />
                      <span>{m.userName}</span>
                    </Table.Cell>
                    <Table.Cell>
                      {m.fullName}
                    </Table.Cell>
                    <Table.Cell style={{paddingLeft: '20' + 'px'}}>
                      {m.isAdmin && <Icon color='green' name='checkmark' />}
                    </Table.Cell>
                    <Table.Cell style={{paddingLeft: '30' + 'px'}}>
                      {m.isDisable && <Icon color='red' name='x icon' />}
                    </Table.Cell>
                    <Table.Cell collapsing >
                    <Popup content='Edit member' trigger={
                      <span>
                        {!m.isDisable && <Icon name='edit' link onClick={() => { this.props.history.push(`/settings/admin/member/${m.uid}/edit`) }} />}
                      </span>
                    }/>
                    <Popup content='Activate member' trigger={
                      <span>
                        {m.isDisable && <Icon name='unlock alternate' link onClick={() => {this.openActivateMember(m.uid)}} />}
                      </span>
                    }/>
                    <Confirm content='Are you sure you want to activate this member?' open={this.state.isOpenActivateMember} onCancel={this.closeActivateMember} onConfirm={this.handleActivateMember} />
                    </Table.Cell>
                  </Table.Row>
          ))}
              </Table.Body>
            </Table>
            { membersQuery.hasMoreData &&
            <Button onClick={() => { membersQuery.loadMoreEntries() }}>Load More Users</Button>
        }
          </div>
        }
      </Container>
    )
  }
}

Members.propTypes = {
}

const MembersQuery = gql`
      query MembersQuery($first: Int, $search: String){
        membersList(first: $first, search: $search) {
          edges {
            cursor
            member {
              uid
              isAdmin
              isDisable
              userName
              fullName
            }
          }
          hasMoreData
        }
      }
`

const MoreMembersQuery = gql`
      query MoreMembersQuery($after: Int){
        membersList(first: $first, after: $after) {
          edges {
            cursor
            member {
              uid
              isAdmin
              userName
              fullName
              isDisable
            }
          }
          hasMoreData
        }
      }
`

const updateActivateMember = gql`
  mutation updateActivateMember($updateActivateMemberChange: UpdateActivateMemberChange!) {
    updateActivateMember(updateActivateMemberChange: $updateActivateMemberChange) {
      hasErrors
      genericError
      member {
        uid
      }
    }
  }
`

const ViewerQuery = gql`
  query viewerQuery {
    viewer {
      member {
        uid
        isAdmin
      }
    }
  }
`
export default compose(
  graphql(updateActivateMember, {
    name: 'updateActivateMember',
    props: ({ updateActivateMember }) => ({
      updateActivateMember: (updateActivateMemberChange) => updateActivateMember({ variables: { updateActivateMemberChange }, refetchQueries: ['memberQuery'] })
    })
  }),
graphql(ViewerQuery, {
  name: 'viewerQuery',
  options: () => ({
    fetchPolicy: 'network-only'
  })
}),
graphql(MembersQuery, {
  options: () => ({
    variables: {
      first: defaultFetchSize
    },
    fetchPolicy: 'network-only'
  }),
  props ({ data: { loading, error, refetch, membersList, fetchMore } }) {
    console.log('loading', loading)
    console.log('error', error)
    let members = membersList
    let cursor
    const cursors = members && members.edges.map((e) => (e.cursor))
    if (cursors && cursors.length > 0) cursor = cursors[cursors.length - 1]
    membersList = members && members.edges.map((e) => (e.member))
    return {
      membersQuery: {
        loading,
        error,
        membersList: membersList,
        hasMoreData: members && members.hasMoreData,
        loadMoreEntries: () => {
          return fetchMore({
            query: MoreMembersQuery,
            variables: {
              first: defaultFetchSize,
              after: cursor
            },
            updateQuery: (previousResult, { fetchMoreResult }) => {
              const newEdges = fetchMoreResult.members.edges

              return {
                membersList: Object.assign({}, fetchMoreResult.members, {
                  edges: [...previousResult.members.edges, ...newEdges],
                  hasMoreData: fetchMoreResult.members.hasMoreData
                })
              }
            }
          })
        },
        update: (searchString) => {
          return fetchMore({
            query: MembersQuery,
            variables: {
              first: defaultFetchSize,
              search: searchString
            },
            updateQuery: (previousResult, { fetchMoreResult }) => {
              console.log('fetchMoreResult', fetchMoreResult)
              const newEdges = fetchMoreResult.members.edges

              return {
                membersList: Object.assign({}, fetchMoreResult.members, {
                  edges: [...newEdges],
                  hasMoreData: fetchMoreResult.members.hasMoreData
                })
              }
            }
          })
        },
        refetch: (searchString) => {
          return fetchMore({
            query: MembersQuery,
            variables: {
              first: defaultFetchSize
            },
            updateQuery: (previousResult, { fetchMoreResult }) => {
              console.log('fetchMoreResult', fetchMoreResult)
              const newEdges = fetchMoreResult.members.edges

              return {
                membersList: Object.assign({}, fetchMoreResult.members, {
                  edges: [...newEdges],
                  hasMoreData: fetchMoreResult.members.hasMoreData
                })
              }
            }
          })
        }
      }
    }
  }
})
)(withError(Members))
