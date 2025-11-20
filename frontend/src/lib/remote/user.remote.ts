import { query } from '$app/server';
import SensorsService from '$lib/services/SensorsService';
import * as v from 'valibot';

const userSensors = ['145799809528704', '39620398887400'];

//TODO: load sensors for the logged in user
export const getUserSensors = query(async () => {
	return SensorsService.getSensorsByFuseIds(userSensors);
});

export const getUserSensorData = query(v.object({ fuseId: v.string(), from: v.date(), to: v.date() }), async ({ fuseId, from, to }) => {
	if (userSensors.includes(fuseId) === false) {
		throw new Error('Access denied to the requested sensor data');
	}
	const result = SensorsService.getSensorData(fuseId, from, to);
	return result;
});
