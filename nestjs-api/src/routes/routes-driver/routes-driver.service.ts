import { Injectable } from '@nestjs/common';
import { RouteDriver } from '../dto/route-driver';
import { PrismaService } from '../../prisma/prisma.service';

@Injectable()
export class RoutesDriverService {
  constructor(private readonly prismaService: PrismaService) {}

  public async processRoute(dto: RouteDriver) {
    const routeDriver = await this.prismaService.routeDriver.upsert({
      include: {
        route: true,
      },
      where: { route_id: dto.route_id },
      create: {
        route_id: dto.route_id,
        points: {
          set: {
            location: {
              lat: dto.lat,
              lng: dto.lng,
            },
          },
        },
      },
      update: {
        points: {
          push: {
            location: {
              lat: dto.lat,
              lng: dto.lng,
            },
          },
        },
      },
    });
    return routeDriver;
  }
}
