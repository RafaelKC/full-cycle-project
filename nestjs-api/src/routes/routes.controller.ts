import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { RoutesService } from './routes.service';
import { CreateRouteDto } from './dto/create-route.dto';

@Controller('routes')
export class RoutesController {
  constructor(private readonly routesService: RoutesService) {}

  @Post()
  public create(@Body() createRouteDto: CreateRouteDto) {
    return this.routesService.create(createRouteDto);
  }

  @Get()
  public findAll() {
    return this.routesService.findAll();
  }

  @Get(':id')
  public findOne(@Param('id') id: string) {
    return this.routesService.findOne(id);
  }
}
