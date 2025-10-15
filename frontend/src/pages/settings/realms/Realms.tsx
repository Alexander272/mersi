import { useState } from 'react'
import { Stack } from '@mui/material'

import { RealmsList } from '@/features/realms/components/RealmsList'
import { RealmForm } from '@/features/realms/components/RealmForm'
import { AccessesTable } from '@/features/accesses/components/AccessesTable'

export default function Realms() {
	const [realm, setRealm] = useState('new')

	const realmHandler = (realm: string) => {
		setRealm(realm)
	}

	return (
		<Stack direction={'row'} spacing={2} height={'100%'}>
			<RealmsList realm={realm} setRealm={realmHandler} />
			<Stack width={'100%'} spacing={3} sx={{ maxHeight: 760, overflowY: 'auto', pt: 1 }}>
				<RealmForm realm={realm} setRealm={realmHandler} />
				{realm != 'new' && <AccessesTable realm={realm} />}
			</Stack>
		</Stack>
	)
}
