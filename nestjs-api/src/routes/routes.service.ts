import { Injectable } from '@nestjs/common';
import { CreateRouteDto } from './dto/create-route.dto';
import { PrismaService } from '../prisma/prisma.service';
import { DirectionsService } from '../maps/diractions/directions.service';

@Injectable()
export class RoutesService {
  constructor(
    private readonly prismaService: PrismaService,
    private readonly directionsService: DirectionsService,
  ) {}

  public async create(createRouteDto: CreateRouteDto) {
    const { available_travel_modes, geocoded_waypoints, routes, request } =
      await this.directionsService.getDirections(
        createRouteDto.source_id,
        createRouteDto.destination_id,
      );

    const legs = routes[0].legs[0];
    return await this.prismaService.route.create({
      data: {
        name: createRouteDto.name,
        source: {
          name: legs.start_address,
          location: legs.start_location,
        },
        destination: {
          name: legs.end_address,
          location: legs.end_location,
        },
        distance: legs.distance.value,
        duration: legs.duration.value,
        directions: JSON.parse(
          JSON.stringify({
            available_travel_modes,
            geocoded_waypoints,
            routes,
            request,
          }),
        ),
      },
    });
  }

  public findAll() {
    return this.prismaService.route.findMany();
  }

  public findOne(id: string) {
    return this.prismaService.route.findUniqueOrThrow({
      where: { id },
    });
  }
}
