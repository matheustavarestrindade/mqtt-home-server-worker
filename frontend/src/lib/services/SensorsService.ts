import { env } from '$env/dynamic/private';

type SensorsByFuseIdResponse = Array<{
	id: number;
	fuse_id: string;
	name: string;

	type: string;
	location: string;
	wifi_strength: number;
	battery_percent: number;

	description: string;
	created_at: string;
	last_seen: string;
}>;

type Sensors = Array<{
	id: number;
	fuse_id: string;
	name: string;
	type: string;
	location: string;
	wifi_strength: number;
	battery_percent: number;
	description: string;
	created_at: Date;
	last_seen: Date;
}>;

type HidroponicManagerSensorDataResponse = {
	payload_version: number;
	temperature: number;
	temperature_severity: number;
	moisture: number;
	moisture_severity: number;
	ph: number;
	ph_severity: number;
	conductivity: number;
	conductivity_severity: number;
	nitrogen: number;
	nitrogen_severity: number;
	phosphorus: number;
	phosphorus_severity: number;
	potassium: number;
	potassium_severity: number;
	isOn: boolean;
	next_toggle_in_seconds: number;
};

type WaterLevelMeterSensorDataResponse = {
	average_water_level_cm: number;
};

type SensorsTypes = 'hidroponic-manager' | 'water-level-meter';

type SensorNotFoundError = {
	error: 'Sensor data not found';
};

export type SensorDataResponse<SensorType> = SensorType extends SensorsTypes
	? SensorType extends 'hidroponic-manager'
		? HidroponicManagerSensorDataResponse[]
		: SensorType extends 'water-level-meter'
			? WaterLevelMeterSensorDataResponse[]
			: never
	: never;

class SensorsService {
	private static readonly URL: string = env.SERVICE_BASE_URL;

	static async getSensorsByFuseIds(fuseIds: string[]): Promise<Sensors> {
		const endpoint = new URL(SensorsService.URL + '/sensors');
		endpoint.searchParams.append('ids', fuseIds.join(','));
		console.log('Fetching sensors for fuse IDs:', endpoint.toString());

		const response = await fetch(endpoint.toString());
		if (!response.ok) {
			throw new Error(`Error fetching sensors: ${response.statusText}`);
		}

		const data: SensorsByFuseIdResponse = await response.json();
		if (!data || data.length === 0) {
			return [];
		}

		console.log('Fetched sensors data:', data);

		return data.map((sensor) => ({
			...sensor,
			created_at: new Date(sensor.created_at),
			last_seen: new Date(sensor.last_seen)
		}));
	}

	static async getSensorData<SensorType>(fuseId: string, from: Date, to: Date): Promise<SensorDataResponse<SensorType> | SensorNotFoundError> {
		const endpoint = new URL(SensorsService.URL + `/sensor/data`);

		endpoint.searchParams.append('fuse_id', fuseId);
		endpoint.searchParams.append('start', from.toISOString());
		endpoint.searchParams.append('end', to.toISOString());
		endpoint.searchParams.append('interval_ms', 1000 * 60 * 5 + ''); // 5 minutes

		const response = await fetch(endpoint.toString());
		if (!response.ok) {
			return { error: 'Sensor data not found' } as SensorNotFoundError;
		}

		const data = await response.json();
		return data;
	}
}

export default SensorsService;
