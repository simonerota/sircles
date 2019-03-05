import React, { PropTypes } from 'react'
import { graphql, compose } from 'react-apollo'
import gql from 'graphql-tag'
import { Link } from 'react-router-dom'
import { Menu, Dropdown, Dimmer, Loader, Icon } from 'semantic-ui-react'
import DayPicker from 'react-day-picker'

import 'react-day-picker/lib/style.css'

import { withError } from '../modules/Error'
import Util from '../modules/Util'

import TimeTravel from './TimeTravel'

class UserMenu extends React.Component {
  componentWillMount () {
    this.resetComponent()
  }

  resetComponent = () => this.setState(
    {
      selectedDay: null
    })

  componentWillReceiveProps (nextProps) {
    const { memberQuery } = nextProps

    if (Util.isQueriesError(memberQuery)) {
      this.props.appError.setError(true)
      return
    }

    // reset selectedDay or it will cause a loop continuously rendering the TimeTravel component that will call history.push etc...
    if (this.props.location !== nextProps.location) {
      this.setState({
        selectedDay: null
      })
    }
  }

  handleDayClick = day => {
    this.setState({
      selectedDay: day
    })
  }

  render () {
    const { memberQuery } = this.props
    const { selectedDay } = this.state
    const timeLine = this.props.match.params.timeLine

    console.log('selectedDay', selectedDay)

    if (selectedDay) {
      return <TimeTravel day={selectedDay} />
    }

    if (memberQuery.error) {
      return null
    }

    if (memberQuery.loading) {
      return (
        <Dimmer active inverted>
          <Loader inverted>Loading</Loader>
        </Dimmer>
      )
    }

    let member
    let rolesMap = {}
    let parents = {}

    // the member may not exist in this timeline
    if (memberQuery.member) {
      member = JSON.parse(JSON.stringify(memberQuery.member))

      // sort circle by depth then by name
      member.circles.sort((a, b) => {
        const d = a.role.depth - b.role.depth
        if (d !== 0) return d
        return a.role.name.localeCompare(b.role.name)
      })

      // show roles grouped by parent
      for (let i = 0; i < member.roles.length; i++) {
        const role = member.roles[i].role
        const parent = role.parent
        if (!rolesMap[parent.uid]) {
          rolesMap[parent.uid] = []
        }
        rolesMap[parent.uid].push(role)
        parents[parent.uid] = parent
      }
    }

    // sort parents by depth
    const parentOrderedKeys = Object.keys(parents).sort((a, b) => (parents[a].depth - parents[b].depth))
    
    // get core circles of the current user
    // and lead links
    let coreCircleUid = 0;
    let sircleLeaderUid = 0;
    for(let i = 0; i < member.roles.length; i++) {
      if(member.roles[i].role.name == "Core Members")
        coreCircleUid = member.roles[i].role.parent.uid
      if(member.roles[i].role.name == "Sircle Leader")
        sircleLeaderUid = member.roles[i].role.parent.uid
    }
    
    return (
      <Menu>
        <Menu.Item name='Organization' as={Link} to={Util.orgChartUrl(null, timeLine)} />
        { member &&
        <Dropdown item scrolling text='My Circles'>
          <Dropdown.Menu>
            { member.circles.map(circle => (
              <Dropdown.Item key={circle.role.uid} as={Link} to={Util.roleUrl(circle.role.uid, timeLine)}>
                {circle.role.name}
                &nbsp;{ circle.role.uid == coreCircleUid && <Icon inverted color='orange' name='selected radio' style={{marginRight: 0 + 'px'}}></Icon>}
                &nbsp;{ circle.role.uid == sircleLeaderUid && <Icon inverted color='blue' name='selected radio' style={{marginRight: 0 + 'px'}}></Icon>}
              </Dropdown.Item>
            ))}
          </Dropdown.Menu>
        </Dropdown>
        }
        { member &&
        <Dropdown item scrolling text='My Roles'>
          <Dropdown.Menu>
            { parentOrderedKeys.map((parentUID) => {
              let roles = rolesMap[parentUID]
              const parent = parents[parentUID]
              // remove core members from roles
              for(let i = 0; i < roles.length; i++) {
                if(roles[i].name === "Core Members") {
                  roles.splice(i, 1)
                }
              }
              const items = roles.map(role => (
                <Dropdown.Item key={role.uid} as={Link} to={Util.roleUrl(role.uid, timeLine)}>
                  {role.name}
                </Dropdown.Item>
              ))
              return (
                [
                  <Dropdown.Header key={parent.uid}>
                    {parent.name}
                    &nbsp;{ parent.uid == coreCircleUid && <Icon inverted color='orange' name='selected radio' style={{marginRight: 0 + 'px'}}></Icon> }
                    &nbsp;{ parent.uid == sircleLeaderUid && <Icon inverted color='blue' name='selected radio' style={{marginRight: 0 + 'px'}}></Icon> }                
                  </Dropdown.Header>,
                  [...items]
                ]
              )
            }
            )}
          </Dropdown.Menu>
        </Dropdown>
        }
        { member &&
        <Menu.Item name='My Tensions' as={Link} to='/tensions' />
        }
        <Menu.Menu position='right' onClick={e => e.stopPropagation()}>
          <Dropdown simple item text='Time Travel' onClick={e => e.stopPropagation()}>
            <Dropdown.Menu>
              <DayPicker onDayClick={this.handleDayClick} />
            </Dropdown.Menu>
          </Dropdown>
        </Menu.Menu>
      </Menu>
    )
  }
}

UserMenu.propTypes = {
  viewer: PropTypes.object.isRequired
}

const MemberQuery = gql`
  query memberQuery($timeLineID: TimeLineID, $uid: ID!) {
    member(timeLineID: $timeLineID, uid: $uid) {
      uid
      isAdmin
      userName
      circles {
        role {
          uid
          depth
          name
        }
        isLeadLink
      }
      roles {
        role {
          uid
          depth
          name
          parent {
            uid
            depth
            name
          }
        }
      }
    }
  }
`

export default compose(
graphql(MemberQuery, {
  name: 'memberQuery',
  options: props => ({
    variables: {
      uid: props.viewer.member.uid,
      timeLineID: props.match.params.timeLine || 0
    },
    fetchPolicy: 'network-only'
  })
})
)(withError(UserMenu))
