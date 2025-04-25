import { Suspense } from 'react'

import { Fallback } from '@/components/Fallback/Fallback'
import { PageBox } from '@/components/PageBox/PageBox'

export default function Home() {
	return (
		<PageBox>
			<Suspense fallback={<Fallback />}>{/* <Table /> */}</Suspense>
		</PageBox>
	)
}
